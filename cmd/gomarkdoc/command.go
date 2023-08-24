package main

import (
	"bytes"
	"container/list"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"hash/fnv"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/format"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
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
	repository                 lang.Repo
	output                     string
	header                     string
	headerFile                 string
	footer                     string
	footerFile                 string
	format                     string
	tags                       []string
	excludeDirs                []string
	excludeLinkAngularBrackets bool
	templateOverrides          map[string]string
	templateFileOverrides      map[string]string
	verbosity                  int
	includeUnexported          bool
	check                      bool
	embed                      bool
	version                    bool
}

// Flags populated by goreleaser
var version = ""

const configFilePrefix = ".gomarkdoc"

func buildCommand() *cobra.Command {
	var opts commandOptions
	var configFile string

	// cobra.OnInitialize(func() { buildConfig(configFile) })

	var command = &cobra.Command{
		Use:   "gomarkdoc [package ...]",
		Short: "generate markdown documentation for golang code",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.version {
				printVersion()
				return nil
			}

			buildConfig(configFile)

			// Load configuration from viper
			opts.includeUnexported = viper.GetBool("includeUnexported")
			opts.output = viper.GetString("output")
			opts.check = viper.GetBool("check")
			opts.embed = viper.GetBool("embed")
			opts.format = viper.GetString("format")
			opts.templateOverrides = viper.GetStringMapString("template")
			opts.templateFileOverrides = viper.GetStringMapString("templateFile")
			opts.header = viper.GetString("header")
			opts.headerFile = viper.GetString("headerFile")
			opts.footer = viper.GetString("footer")
			opts.footerFile = viper.GetString("footerFile")
			opts.tags = viper.GetStringSlice("tags")
			opts.excludeDirs = viper.GetStringSlice("excludeDirs")
			opts.excludeLinkAngularBrackets = viper.GetBool("excludeLinkAngularBrackets")
			opts.repository.Remote = viper.GetString("repository.url")
			opts.repository.DefaultBranch = viper.GetString("repository.defaultBranch")
			opts.repository.PathFromRoot = viper.GetString("repository.path")

			if opts.check && opts.output == "" {
				return errors.New("gomarkdoc: check mode cannot be run without an output set")
			}

			if len(args) == 0 {
				// Default to current directory
				args = []string{"."}
			}

			return runCommand(args, opts)
		},
	}

	command.Flags().StringVar(
		&configFile,
		"config",
		"",
		fmt.Sprintf("File from which to load configuration (default: %s.yml)", configFilePrefix),
	)
	command.Flags().BoolVarP(
		&opts.includeUnexported,
		"include-unexported",
		"u",
		false,
		"Output documentation for unexported symbols, methods and fields in addition to exported ones.",
	)
	command.Flags().StringVarP(
		&opts.output,
		"output",
		"o",
		"",
		"File or pattern specifying where to write documentation output. Defaults to printing to stdout.",
	)
	command.Flags().BoolVarP(
		&opts.check,
		"check",
		"c",
		false,
		"Check the output to see if it matches the generated documentation. --output must be specified to use this.",
	)
	command.Flags().BoolVarP(
		&opts.embed,
		"embed",
		"e",
		false,
		"Embed documentation into existing markdown files if available, otherwise append to file.",
	)
	command.Flags().StringVarP(
		&opts.format,
		"format",
		"f",
		"github",
		"Format to use for writing output data. Valid options: github (default), azure-devops, plain",
	)
	command.Flags().StringToStringVarP(
		&opts.templateOverrides,
		"template",
		"t",
		map[string]string{},
		"Custom template string to use for the provided template name instead of the default template.",
	)
	command.Flags().StringToStringVar(
		&opts.templateFileOverrides,
		"template-file",
		map[string]string{},
		"Custom template file to use for the provided template name instead of the default template.",
	)
	command.Flags().StringVar(
		&opts.header,
		"header",
		"",
		"Additional content to inject at the beginning of each output file.",
	)
	command.Flags().StringVar(
		&opts.headerFile,
		"header-file",
		"",
		"File containing additional content to inject at the beginning of each output file.",
	)
	command.Flags().StringVar(
		&opts.footer,
		"footer",
		"",
		"Additional content to inject at the end of each output file.",
	)
	command.Flags().StringVar(
		&opts.footerFile,
		"footer-file",
		"",
		"File containing additional content to inject at the end of each output file.",
	)
	command.Flags().StringSliceVar(
		&opts.tags,
		"tags",
		defaultTags(),
		"Set of build tags to apply when choosing which files to include for documentation generation.",
	)
	command.Flags().StringSliceVar(
		&opts.excludeDirs,
		"exclude-dirs",
		nil,
		"List of package directories to ignore when producing documentation.",
	)
	command.Flags().BoolVarP(
		&opts.excludeLinkAngularBrackets,
		"exclude-link-angular-brackets",
		"",
		false,
		"Exclude the angular brackets [](<>) from links in the output. Works with the github format only.",
	)
	command.Flags().CountVarP(
		&opts.verbosity,
		"verbose",
		"v",
		"Log additional output from the execution of the command. Can be chained for additional verbosity.",
	)
	command.Flags().StringVar(
		&opts.repository.Remote,
		"repository.url",
		"",
		"Manual override for the git repository URL used in place of automatic detection.",
	)
	command.Flags().StringVar(
		&opts.repository.DefaultBranch,
		"repository.default-branch",
		"",
		"Manual override for the git repository URL used in place of automatic detection.",
	)
	command.Flags().StringVar(
		&opts.repository.PathFromRoot,
		"repository.path",
		"",
		"Manual override for the path from the root of the git repository used in place of automatic detection.",
	)
	command.Flags().BoolVar(
		&opts.version,
		"version",
		false,
		"Print the version.",
	)

	// We ignore the errors here because they only happen if the specified flag doesn't exist
	_ = viper.BindPFlag("includeUnexported", command.Flags().Lookup("include-unexported"))
	_ = viper.BindPFlag("output", command.Flags().Lookup("output"))
	_ = viper.BindPFlag("check", command.Flags().Lookup("check"))
	_ = viper.BindPFlag("embed", command.Flags().Lookup("embed"))
	_ = viper.BindPFlag("format", command.Flags().Lookup("format"))
	_ = viper.BindPFlag("template", command.Flags().Lookup("template"))
	_ = viper.BindPFlag("templateFile", command.Flags().Lookup("template-file"))
	_ = viper.BindPFlag("header", command.Flags().Lookup("header"))
	_ = viper.BindPFlag("headerFile", command.Flags().Lookup("header-file"))
	_ = viper.BindPFlag("footer", command.Flags().Lookup("footer"))
	_ = viper.BindPFlag("footerFile", command.Flags().Lookup("footer-file"))
	_ = viper.BindPFlag("tags", command.Flags().Lookup("tags"))
	_ = viper.BindPFlag("excludeDirs", command.Flags().Lookup("exclude-dirs"))
	_ = viper.BindPFlag("excludeLinkAngularBrackets", command.Flags().Lookup("exclude-link-angular-brackets"))
	_ = viper.BindPFlag("repository.url", command.Flags().Lookup("repository.url"))
	_ = viper.BindPFlag("repository.defaultBranch", command.Flags().Lookup("repository.default-branch"))
	_ = viper.BindPFlag("repository.path", command.Flags().Lookup("repository.path"))

	return command
}

func defaultTags() []string {
	f, ok := os.LookupEnv("GOFLAGS")
	if !ok {
		return nil
	}

	fs := flag.NewFlagSet("goflags", flag.ContinueOnError)
	tags := fs.String("tags", "", "")

	if err := fs.Parse(strings.Fields(f)); err != nil {
		return nil
	}

	if tags == nil {
		return nil
	}

	return strings.Split(*tags, ",")
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

	excluded := getSpecs(opts.excludeDirs...)
	if err := validateExcludes(excluded); err != nil {
		return err
	}

	specs = removeExcludes(specs, excluded)

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
		f = format.NewGitHubFlavoredMarkdown(opts.excludeLinkAngularBrackets)
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
		log := logger.New(getLogLevel(opts.verbosity), logger.WithField("dir", spec.Dir))

		buildPkg, err := getBuildPackage(spec.ImportPath, opts.tags)
		if err != nil {
			log.Debugf("unable to load package in directory: %s", err)
			// We don't care if a wildcard path produces nothing
			if spec.isWildcard {
				continue
			}

			return err
		}

		var pkgOpts []lang.PackageOption
		pkgOpts = append(pkgOpts, lang.PackageWithRepositoryOverrides(&opts.repository))

		if opts.includeUnexported {
			pkgOpts = append(pkgOpts, lang.PackageWithUnexportedIncluded())
		}

		pkg, err := lang.NewPackageFromBuild(log, buildPkg, pkgOpts...)
		if err != nil {
			return err
		}

		spec.pkg = pkg
	}

	return nil
}

func getBuildPackage(path string, tags []string) (*build.Package, error) {
	ctx := build.Default
	ctx.BuildTags = tags

	if isLocalPath(path) {
		pkg, err := ctx.ImportDir(path, build.ImportComment)
		if err != nil {
			return nil, fmt.Errorf("gomarkdoc: invalid package in directory: %s", path)
		}

		return pkg, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pkg, err := ctx.Import(path, wd, build.ImportComment)
	if err != nil {
		return nil, fmt.Errorf("gomarkdoc: invalid package at import path: %s", path)
	}

	return pkg, nil
}

func getSpecs(paths ...string) []*PackageSpec {
	var expanded []*PackageSpec
	for _, path := range paths {
		// Ensure that the path we're working with is normalized for the OS
		// we're using (i.e. "\" for windows, "/" for everything else)
		path = filepath.FromSlash(path)

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

// validateExcludes checks that the exclude dirs are all directories, not
// packages.
func validateExcludes(specs []*PackageSpec) error {
	for _, s := range specs {
		if !s.isLocal {
			return fmt.Errorf("gomarkdoc: invalid directory specified as an exclude directory: %s", s.ImportPath)
		}
	}

	return nil
}

// removeExcludes removes any package specs that were specified as excluded.
func removeExcludes(specs []*PackageSpec, excludes []*PackageSpec) []*PackageSpec {
	out := make([]*PackageSpec, 0, len(specs))
	for _, s := range specs {
		var exclude bool
		for _, e := range excludes {
			if !s.isLocal || !e.isLocal {
				continue
			}

			if r, err := filepath.Rel(s.Dir, e.Dir); err == nil && r == "." {
				exclude = true
				break
			}
		}

		if !exclude {
			out = append(out, s)
		}
	}

	return out
}

const (
	cwdPathPrefix    = "." + string(os.PathSeparator)
	parentPathPrefix = ".." + string(os.PathSeparator)
)

func isLocalPath(path string) bool {
	return strings.HasPrefix(path, ".") || strings.HasPrefix(path, parentPathPrefix) || filepath.IsAbs(path)
}

func compare(r1, r2 io.Reader) (bool, error) {
	r1Hash := fnv.New128()
	if _, err := io.Copy(r1Hash, r1); err != nil {
		return false, fmt.Errorf("gomarkdoc: failed when checking documentation: %w", err)
	}

	r2Hash := fnv.New128()
	if _, err := io.Copy(r2Hash, r2); err != nil {
		return false, fmt.Errorf("gomarkdoc: failed when checking documentation: %w", err)
	}

	return bytes.Equal(r1Hash.Sum(nil), r2Hash.Sum(nil)), nil
}

func getLogLevel(verbosity int) logger.Level {
	switch verbosity {
	case 0:
		return logger.WarnLevel
	case 1:
		return logger.InfoLevel
	case 2:
		return logger.DebugLevel
	default:
		return logger.DebugLevel
	}
}

func printVersion() {
	if version != "" {
		fmt.Println(version)
		return
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		fmt.Println(info.Main.Version)
	} else {
		fmt.Println("<unknown>")
	}
}
