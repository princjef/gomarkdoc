package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
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

	for fileName, pkgs := range filePkgs {
		file := lang.NewFile(header, footer, pkgs)

		text, err := out.File(file)
		if err != nil {
			return err
		}

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
				return err
			}
		default:
			if err := writeFile(fileName, text); err != nil {
				return fmt.Errorf("failed to write output file %s: %w", fileName, err)
			}
		}
	}

	return nil
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

	f, err := os.Open(path)
	if err != nil {
		if err == os.ErrNotExist {
			return checkErr
		}

		return fmt.Errorf("failed to open file %s for checking: %w", path, err)
	}

	defer f.Close()

	match, err := compare(b, f)
	if err != nil {
		return fmt.Errorf("failure while attempting to check contents of %s: %w", path, err)
	}

	if !match {
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
