package main

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"go/build"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/format"
	"github.com/princjef/gomarkdoc/lang"
)

// PackageSpec defines the data available to the --output option's template.
// Information is recomputed for each package generated.
type PackageSpec struct {
	// Dir holds the local path where the package is located. If the package is
	// a remote package, this will always be ".".
	Dir string

	// ImportPath holds a representation of the package that should be unique
	// for most purposes. If a package is on the filesystem, this is equivalent
	// to the value of Dir. For remote packages, this holds the string used to
	// import that package in code (e.g. "encoding/json").
	ImportPath string
	isWildcard bool
	isLocal    bool
	outputFile string
	pkg        *lang.Package
}

type commandOptions struct {
	includeUnexported     bool
	output                string
	check                 bool
	templateOverrides     map[string]string
	templateFileOverrides map[string]string
	header                string
	headerFile            string
	footer                string
	footerFile            string
	format                string
}

const configFilePrefix = ".gomarkdoc"

func buildCommand() *cobra.Command {
	var opts commandOptions
	var configFile string

	cobra.OnInitialize(func() { buildConfig(configFile) })

	var command = &cobra.Command{
		Use:   "gomarkdoc [package ...]",
		Short: "generate markdown documentation for golang code",
		Run: func(cmd *cobra.Command, args []string) {
			// Load configuration from viper
			opts.includeUnexported = viper.GetBool("includeUnexported")
			opts.output = viper.GetString("output")
			opts.check = viper.GetBool("check")
			opts.format = viper.GetString("format")
			opts.templateOverrides = viper.GetStringMapString("template")
			opts.templateFileOverrides = viper.GetStringMapString("template-file")
			opts.header = viper.GetString("header")
			opts.headerFile = viper.GetString("headerFile")
			opts.footer = viper.GetString("footer")
			opts.footerFile = viper.GetString("footerFile")

			if opts.check && opts.output == "" {
				log.Fatal("check mode cannot be run without an output set")
			}

			if len(args) == 0 {
				// Default to current directory
				args = []string{"."}
			}

			if err := runCommand(args, opts); err != nil {
				log.Fatal(err)
			}
		},
	}

	command.Flags().StringVar(&configFile, "config", "", fmt.Sprintf("File from which to load configuration (default: %s.yml)", configFilePrefix))
	command.Flags().BoolVarP(&opts.includeUnexported, "include-unexported", "u", false, "Output documentation for unexported symbols, methods and fields in addition to exported ones.")
	command.Flags().StringVarP(&opts.output, "output", "o", "", "File or pattern specifying where to write documentation output. Defaults to printing to stdout.")
	command.Flags().BoolVarP(&opts.check, "check", "c", false, "Check the output to see if it matches the generated documentation. --output must be specified to use this option.")
	command.Flags().StringVarP(&opts.format, "format", "f", "github", "Format to use for writing output data. Valid options: github (default), azure-devops, plain")
	command.Flags().StringToStringVarP(&opts.templateOverrides, "template", "t", map[string]string{}, "Custom template string to use for the provided template name instead of the default template.")
	command.Flags().StringToStringVar(&opts.templateFileOverrides, "template-file", map[string]string{}, "Custom template file to use for the provided template name instead of the default template.")
	command.Flags().StringVar(&opts.header, "header", "", "Additional content to inject at the beginning of each output file.")
	command.Flags().StringVar(&opts.headerFile, "header-file", "", "File containing additional content to inject at the beginning of each output file.")
	command.Flags().StringVar(&opts.footer, "footer", "", "Additional content to inject at the end of each output file.")
	command.Flags().StringVar(&opts.footerFile, "footer-file", "", "File containing additional content to inject at the end of each output file.")

	viper.BindPFlag("includeUnexported", command.Flags().Lookup("include-unexported"))
	viper.BindPFlag("output", command.Flags().Lookup("output"))
	viper.BindPFlag("check", command.Flags().Lookup("check"))
	viper.BindPFlag("format", command.Flags().Lookup("format"))
	viper.BindPFlag("template", command.Flags().Lookup("template"))
	viper.BindPFlag("templateFile", command.Flags().Lookup("template-file"))
	viper.BindPFlag("header", command.Flags().Lookup("header"))
	viper.BindPFlag("headerFile", command.Flags().Lookup("header-file"))
	viper.BindPFlag("footer", command.Flags().Lookup("footer"))
	viper.BindPFlag("footerFile", command.Flags().Lookup("footer-file"))

	return command
}

func buildConfig(configFile string) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(configFilePrefix)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// TODO: better handling
			fmt.Println(err)
		}
	}
}

func runCommand(paths []string, opts commandOptions) error {
	outputTmpl, err := template.New("output").Parse(opts.output)
	if err != nil {
		return fmt.Errorf("gomarkdoc: invalid output template: %w", err)
	}

	specs := getSpecs(paths...)

	if err := resolveOutput(specs, outputTmpl); err != nil {
		return err
	}

	if err := loadPackages(specs, opts); err != nil {
		return err
	}

	return writeOutput(specs, opts)
}

func resolveOutput(specs []*PackageSpec, outputTmpl *template.Template) error {
	for _, spec := range specs {
		var outputFile strings.Builder
		if err := outputTmpl.Execute(&outputFile, spec); err != nil {
			return err
		}

		outputStr := outputFile.String()
		if outputStr == "" {
			// Preserve empty values
			spec.outputFile = ""
		} else {
			// Clean up other values
			spec.outputFile = filepath.Clean(outputFile.String())
		}
	}

	return nil
}

func resolveOverrides(opts commandOptions) ([]gomarkdoc.RendererOption, error) {
	var overrides []gomarkdoc.RendererOption

	// Content overrides take precedence over file overrides
	for name, s := range opts.templateOverrides {
		overrides = append(overrides, gomarkdoc.WithTemplateOverride(name, s))
	}

	for name, f := range opts.templateFileOverrides {
		// File overrides get applied only if there isn't already a content
		// override.
		if _, ok := opts.templateOverrides[name]; ok {
			continue
		}

		b, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("gomarkdoc: couldn't resolve template for %s: %w", name, err)
		}

		overrides = append(overrides, gomarkdoc.WithTemplateOverride(name, string(b)))
	}

	var f format.Format
	switch opts.format {
	case "github":
		f = &format.GitHubFlavoredMarkdown{}
	case "azure-devops":
		f = &format.AzureDevOpsMarkdown{}
	case "plain":
		f = &format.PlainMarkdown{}
	default:
		return nil, fmt.Errorf("gomarkdoc: invalid format: %s", opts.format)
	}

	overrides = append(overrides, gomarkdoc.WithFormat(f))

	return overrides, nil
}

func resolveHeader(opts commandOptions) (string, error) {
	if opts.header != "" {
		return opts.header, nil
	}

	if opts.headerFile != "" {
		b, err := ioutil.ReadFile(opts.headerFile)
		if err != nil {
			return "", fmt.Errorf("gomarkdoc: couldn't resolve header file: %w", err)
		}

		return string(b), nil
	}

	return "", nil
}

func resolveFooter(opts commandOptions) (string, error) {
	if opts.footer != "" {
		return opts.footer, nil
	}

	if opts.footerFile != "" {
		b, err := ioutil.ReadFile(opts.footerFile)
		if err != nil {
			return "", fmt.Errorf("gomarkdoc: couldn't resolve footer file: %w", err)
		}

		return string(b), nil
	}

	return "", nil
}

func loadPackages(specs []*PackageSpec, opts commandOptions) error {
	for _, spec := range specs {
		buildPkg, err := getBuildPackage(spec.ImportPath)
		if err != nil {
			// We don't care if a wildcard path produces nothing
			if spec.isWildcard {
				continue
			}

			return err
		}

		var pkgOpts []lang.PackageOption
		if opts.includeUnexported {
			pkgOpts = append(pkgOpts, lang.PackageWithUnexportedIncluded())
		}

		pkg, err := lang.NewPackageFromBuild(buildPkg, pkgOpts...)
		if err != nil {
			return err
		}

		spec.pkg = pkg
	}

	return nil
}

func writeOutput(specs []*PackageSpec, opts commandOptions) error {
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

		if fileName == "" {
			fmt.Fprint(os.Stdout, text)
		} else if opts.check {
			var b bytes.Buffer
			fmt.Fprint(&b, text)
			if err := checkFile(&b, fileName); err != nil {
				return err
			}
		} else {
			if err := ioutil.WriteFile(fileName, []byte(text), 0755); err != nil {
				return fmt.Errorf("Failed to write output file %s: %w", fileName, err)
			}
		}
	}

	return nil
}

func checkFile(b *bytes.Buffer, path string) error {
	checkErr := errors.New("output does not match current files. Did you forget to run gomarkdoc?")

	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		if err == os.ErrNotExist {
			return checkErr
		} else {
			return fmt.Errorf("Failed to open file %s for checking: %w", path, err)
		}
	}

	match, err := compare(b, f)
	if err != nil {
		return fmt.Errorf("Failure while attempting to check contents of %s: %w", path, err)
	}

	if !match {
		return checkErr
	}

	return nil
}

func getBuildPackage(path string) (*build.Package, error) {
	if isLocalPath(path) {
		pkg, err := build.ImportDir(path, build.ImportComment)
		if err != nil {
			return nil, fmt.Errorf("gomarkdoc: invalid package in directory: %s", path)
		}

		return pkg, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pkg, err := build.Import(path, wd, build.ImportComment)
	if err != nil {
		return nil, fmt.Errorf("gomarkdoc: invalid package at import path: %s", path)
	}

	return pkg, nil
}

func getSpecs(paths ...string) []*PackageSpec {
	var expanded []*PackageSpec
	for _, path := range paths {
		// Not a recursive path
		if !strings.HasSuffix(path, fmt.Sprintf("%s...", string(os.PathSeparator))) {
			isLocal := isLocalPath(path)
			var dir string
			if isLocal {
				dir = path
			} else {
				dir = "."
			}
			expanded = append(expanded, &PackageSpec{
				Dir:        dir,
				ImportPath: path,
				isWildcard: false,
				isLocal:    isLocal,
			})
			continue
		}

		// Remove the recursive marker so we can work with the path
		trimmedPath := path[0 : len(path)-3]

		// Not a file path. Add the original path back to the list so as to not
		// mislead someone into thinking we're processing the recursive path
		if !isLocalPath(trimmedPath) {
			expanded = append(expanded, &PackageSpec{
				Dir:        ".",
				ImportPath: path,
				isWildcard: false,
				isLocal:    false,
			})
			continue
		}

		expanded = append(expanded, &PackageSpec{
			Dir:        trimmedPath,
			ImportPath: trimmedPath,
			isWildcard: true,
			isLocal:    true,
		})

		queue := list.New()
		queue.PushBack(trimmedPath)
		for e := queue.Front(); e != nil; e = e.Next() {
			prev := e.Prev()
			if prev != nil {
				queue.Remove(prev)
			}

			p := e.Value.(string)

			files, err := ioutil.ReadDir(p)
			if err != nil {
				// If we couldn't read the folder, there are no directories that
				// we're going to find beneath it
				continue
			}

			for _, f := range files {
				if isIgnoredDir(f.Name()) {
					continue
				}

				if f.IsDir() {
					subPath := filepath.Join(p, f.Name())

					// Some local paths have their prefixes stripped by Join().
					// If the path is no longer a local path, add the current
					// working directory.
					if !isLocalPath(subPath) {
						subPath = fmt.Sprintf("%s%s", cwdPathPrefix, subPath)
					}

					expanded = append(expanded, &PackageSpec{
						Dir:        subPath,
						ImportPath: subPath,
						isWildcard: true,
						isLocal:    true,
					})
					queue.PushBack(subPath)
				}
			}
		}
	}

	return expanded
}

var ignoredDirs = []string{".git"}

// isIgnoredDir identifies if the dir is one we want to intentionally ignore.
func isIgnoredDir(dirname string) bool {
	for _, ignored := range ignoredDirs {
		if ignored == dirname {
			return true
		}
	}

	return false
}

const (
	cwdPathPrefix    = "." + string(os.PathSeparator)
	parentPathPrefix = ".." + string(os.PathSeparator)
)

func isLocalPath(path string) bool {
	return strings.HasPrefix(path, cwdPathPrefix) || strings.HasPrefix(path, parentPathPrefix) || filepath.IsAbs(path)
}

func compare(r1, r2 io.Reader) (bool, error) {
	b1 := make([]byte, 1024)
	b2 := make([]byte, 1024)

	var count1 int
	var count2 int

	var offset1 int
	var offset2 int

	start := true

	for start || count1 > 0 || count2 > 0 {
		var err error
		// Phase 1: read data if necessary
		if count1 == 0 {
			count1, err = r1.Read(b1)
			if err != nil {
				if err != io.EOF {
					return false, fmt.Errorf("gomarkdoc: failed when checking documentation: %w", err)
				}

				// If the other buffer has more data and we're done, they're not
				// equal
				if count1 == 0 && count2 > 0 {
					return false, nil
				}
			}

			offset1 = 0
		}

		if count2 == 0 {
			count2, err = r2.Read(b2)
			if err != nil {
				if err != io.EOF {
					return false, fmt.Errorf("gomarkdoc: failed when checking documentation: %w", err)
				}

				// If the other buffer has more data and we're done, they're not
				// equal
				if count2 == 0 && count1 > 0 {
					return false, nil
				}
			}

			offset2 = 0
		}

		// Phase 2: compare buffers
		var bytesToRead int
		if count1 < count2 {
			bytesToRead = count1
		} else {
			bytesToRead = count2
		}

		for i := 0; i < bytesToRead; i++ {
			if b1[offset1+i] != b2[offset2+i] {
				return false, nil
			}
		}

		// Phase 3: update counters
		count1 -= bytesToRead
		count2 -= bytesToRead
		offset1 += bytesToRead
		offset2 += bytesToRead

		start = false
	}

	return true, nil
}
