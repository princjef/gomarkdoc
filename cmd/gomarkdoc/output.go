package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
	"github.com/princjef/termdiff"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func writeOutput(specs []*PackageSpec, opts commandOptions) error {
	log := logger.New(getLogLevel(opts.verbosity))

	overrides, err := resolveOverrides(opts)
	if err != nil {
		return err
	}

	out, err := gomarkdoc.NewRenderer(overrides...)
	if err != nil {
		return err
	}

	header, err := resolveHeader(opts)
	if err != nil {
		return err
	}

	footer, err := resolveFooter(opts)
	if err != nil {
		return err
	}

	filePkgs := make(map[string][]*lang.Package)

	for _, spec := range specs {
		if spec.pkg == nil {
			continue
		}

		filePkgs[spec.outputFile] = append(filePkgs[spec.outputFile], spec.pkg)
	}

	var checkErr error
	for fileName, pkgs := range filePkgs {
		file := lang.NewFile(header, footer, pkgs)

		text, err := out.File(file)
		if err != nil {
			return err
		}

		checkErr, err = handleFile(log, fileName, text, opts)
		if err != nil {
			return err
		}
	}

	if checkErr != nil {
		return checkErr
	}

	return nil
}

func handleFile(log logger.Logger, fileName string, text string, opts commandOptions) (error, error) {
	if opts.embed && fileName != "" {
		text = embedContents(log, fileName, text)
	}

	switch {
	case fileName == "":
		fmt.Fprint(os.Stdout, text)
	case opts.check:
		var b bytes.Buffer
		fmt.Fprint(&b, text)
		if err := checkFile(&b, fileName); err != nil {
			return err, nil
		}
	default:
		if err := writeFile(fileName, text); err != nil {
			return nil, fmt.Errorf("failed to write output file %s: %w", fileName, err)
		}
	}
	return nil, nil
}

func writeFile(fileName string, text string) error {
	folder := filepath.Dir(fileName)

	if folder != "" {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return fmt.Errorf("failed to create folder %s: %w", folder, err)
		}
	}

	if err := ioutil.WriteFile(fileName, []byte(text), 0664); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fileName, err)
	}

	return nil
}

func checkFile(b *bytes.Buffer, path string) error {
	checkErr := errors.New("output does not match current files. Did you forget to run gomarkdoc?")

	fileContents, err := os.ReadFile(path)
	if err == os.ErrNotExist {
		fileContents = []byte{}
	} else if err != nil {
		return fmt.Errorf("failed to open file %s for checking: %w", path, err)
	}

	differ := diffmatchpatch.New()
	diff := differ.DiffBisect(b.String(), string(fileContents), time.Now().Add(time.Second))

	// Remove equal diffs
	var filtered = make([]diffmatchpatch.Diff, 0, len(diff))
	for _, d := range diff {
		if d.Type == diffmatchpatch.DiffEqual {
			continue
		}

		filtered = append(filtered, d)
	}

	if len(filtered) != 0 {
		diffs := termdiff.DiffsFromDiffMatchPatch(diff)
		fmt.Fprintln(os.Stderr)
		termdiff.Fprint(
			os.Stderr,
			path,
			diffs,
			termdiff.WithBeforeText("(expected)"),
			termdiff.WithAfterText("(actual)"),
		)
		return checkErr
	}

	return nil
}

var (
	embedStandaloneRegex = regexp.MustCompile(`(?m:^ *)<!--\s*gomarkdoc:embed\s*-->(?m:\s*?$)`)
	embedStartRegex      = regexp.MustCompile(
		`(?m:^ *)<!--\s*gomarkdoc:embed:start\s*-->(?s:.*?)<!--\s*gomarkdoc:embed:end\s*-->(?m:\s*?$)`,
	)
)

func embedContents(log logger.Logger, fileName string, text string) string {
	embedText := fmt.Sprintf("<!-- gomarkdoc:embed:start -->\n\n%s\n\n<!-- gomarkdoc:embed:end -->", text)

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Debugf("unable to find output file %s for embedding. Creating a new file instead", fileName)
		return embedText
	}

	var replacements int
	data = embedStandaloneRegex.ReplaceAllFunc(data, func(_ []byte) []byte {
		replacements++
		return []byte(embedText)
	})

	data = embedStartRegex.ReplaceAllFunc(data, func(_ []byte) []byte {
		replacements++
		return []byte(embedText)
	})

	if replacements == 0 {
		log.Debugf("no embed markers found. Appending documentation to the end of the file instead")
		return fmt.Sprintf("%s\n\n%s", string(data), text)
	}

	return string(data)
}
