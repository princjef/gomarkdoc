// Package gomarkdoc formats documentation for one or more packages as markdown
// for usage outside of the main https://pkg.go.dev site. It supports custom
// templates for tweaking representation of documentation at fine-grained
// levels, exporting both exported and unexported symbols, and custom formatters
// for different backends.
//
// Command Line Usage
//
// If you want to use this package as a command-line tool, you can install the
// command by running:
//
//	go get -u github.com/princjef/gomarkdoc/cmd/gomarkdoc
//
// The command line tool supports configuration for all of the features of the
// importable package:
//
//	$ gomarkdoc --help
//	generate markdown documentation for golang code
//
// 	Usage:
// 	  gomarkdoc [flags] [package ...]
//
// 	Flags:
// 	  -c, --check                              Check the output to see if it matches the generated documentation. --output must be specified to use this.
// 	      --config string                      File from which to load configuration (default: .gomarkdoc.yml)
// 	      --footer string                      Additional content to inject at the end of each output file.
// 	      --footer-file string                 File containing additional content to inject at the end of each output file.
// 	  -f, --format string                      Format to use for writing output data. Valid options: github (default), azure-devops, plain (default "github")
// 	      --header string                      Additional content to inject at the beginning of each output file.
// 	      --header-file string                 File containing additional content to inject at the beginning of each output file.
// 	  -h, --help                               help for gomarkdoc
// 	  -u, --include-unexported                 Output documentation for unexported symbols, methods and fields in addition to exported ones.
// 	  -o, --output string                      File or pattern specifying where to write documentation output. Defaults to printing to stdout.
// 	      --repository.default-branch string   Manual override for the git repository URL used in place of automatic detection.
// 	      --repository.path string             Manual override for the path from the root of the git repository used in place of automatic detection.
// 	      --repository.url string              Manual override for the git repository URL used in place of automatic detection.
// 	  -t, --template stringToString            Custom template string to use for the provided template name instead of the default template. (default [])
// 	      --template-file stringToString       Custom template file to use for the provided template name instead of the default template. (default [])
// 	  -v, --verbose count                      Log additional output from the execution of the command. Can be chained for additional verbosity.
// 	      --version                            Print the version.
//
// The gomarkdoc command processes each of the provided packages, generating
// documentation for the package in markdown format and writing it to console.
// For example, if you have a package in your current directory and want to
// send it to a documentation markdown file, you might do something like this:
//
//	gomarkdoc . > doc.md
//
// Package Specifiers
//
// The gomarkdoc tool supports generating documentation for both local packages
// and remote ones. To specify a local package, start the name of the package
// with a period (.) or specify an absolute path on the filesystem. All other
// package signifiers are assumed to be remote packages. You may specify both
// local and remote packages in the same command invocation as separate
// arguments.
//
// Output Redirection
//
// If you want to redirect output for each processed package to a file, you can
// alternatively provide the --output/-o option, which accepts a template
// specifying how to generate the path of the output file. A common usage of
// this option is when generating README documentation for a package with
// subpackages (which are supported via the ... signifier available in other
// tools):
//
//	gomarkdoc --output '{{.Dir}}/README.md' ./...
//
// You can see all of the data available to the output template in the
// PackageSpec struct in the github.com/princjef/gomarkdoc/cmd/gomarkdoc
// package.
//
// Template Overrides
//
// The documentation information that is output is formatted using a series of
// text templates for the various components of the overall documentation which
// get generated. Higher level templates contain lower level templates, but
// any template may be replaced with an override template using the
// --template/-t option. The full list of templates that may be overridden are:
//
//	- file:    generates documentation for a file containing one or more
//	           packages, depending on how the tool is configured. This is the
//	           root template for documentation generation.
//
//	- package: generates documentation for an entire package.
//
//	- type:    generates documentation for a single type declaration, as well
//	           as any related functions/methods.
//
//	- func:    generates documentation for a single function or method. It may
//	           be referenced from within a type, or directly in the package,
//	           depending on nesting.
//
//	- value:   generates documentation for a single variable or constant
//	           declaration block within a package.
//
//	- index:   generates an index of symbols within a package, similar to what
//	           is seen for godoc.org. The index links to types, funcs,
//	           variables, and constants generated by other templates, so it may
//	           need to be overridden as well if any of those templates are
//	           changed in a material way.
//
//	- example: generates documentation for a single example for a package or
//	           one of its symbols. The example is generated alongside whichever
//	           symbol it represents, based on the standard naming conventions
//	           outlined in https://blog.golang.org/examples#TOC_4.
//
//	- doc:     generates the freeform documentation block for any of the above
//	           structures that can contain a documentation section.
//
// Overriding with the -t option uses a key-vaule pair mapping a template name
// to the file containing the contents of the override template to use.
// Specified template files must exist:
//
//	gomarkdoc -t package=custom-package.gotxt -t doc=custom-doc.gotxt .
//
// Additional Options
//
// As with the godoc tool itself, only exported symbols will be shown in
// documentation. This can be expanded to include all symbols in a package by
// adding the --include-unexported/-u flag.
//
//	gomarkdoc -u . > README.md
//
// You can also run gomarkdoc in a verification mode with the --check/-c flag.
// This is particularly useful for continuous integration when you want to make
// sure that a commit correctly updated the generated documentation. This flag
// is only supported when the --output/-o flag is specified, as the file
// provided there is what the tool is checking:
//
//	gomarkdoc -o README.md -c .
//
// If you're experiencing difficulty with gomarkdoc or just want to get more
// information about how it's executing underneath, you can add -v to show more
// logs. This can be chained a second time to show even more verbose logs:
//
//	gomarkdoc -vv -o README.md .
//
// Some features of gomarkdoc rely on being able to detect information from the
// git repository containing the project. Since individual local git
// repositories may be configured differently from person to person, you may
// want to manually specify the information for the repository to remove any
// inconsistencies. This can be achieved with the --repository.url,
// --repository.default-branch and --repository.path options. For example, this
// repository would be configured with:
//
//	gomarkdoc --repository.url "https://github.com/princjef/gomarkdoc" --repository.defaultBranch master --repository.path / -o README.md .
//
// Configuring via File
//
// If you want to reuse configuration options across multiple invocations, you
// can specify a file in the folder where you invoke gomarkdoc containing
// configuration information that you would otherwise provide on the command
// line. This file may be a JSON, TOML, YAML, HCL, env, or Java properties
// file, but the name is expected to start with .gomarkdoc (e.g.
// .gomarkdoc.yml).
//
// All configuration options are available with the camel-cased form of their
// long name (e.g. --include-unexported becomes includeUnexported). Template
// overrides are specified as a map, rather than a set of key-value pairs
// separated by =. Options provided on the command line override those provided
// in the configuration file if an option is present in both.
//
// Programmatic Usage
//
// While most users will find the command line utility sufficient for their
// needs, this package may also be used programmatically by installing it
// directly, rather than its command subpackage. The programmatic usage
// provides more flexibility when selecting what packages to work with and what
// components to generate documentation for.
//
// A common usage will look something like this:
//
//	package main
//
//	import (
//		"go/build"
//		"fmt"
//		"os"
//
//		"github.com/princjef/gomarkdoc"
//		"github.com/princjef/gomarkdoc/lang"
//	)
//
//	func main() {
//		// Create a renderer to output data
//		out, err := gomarkdoc.NewRenderer()
//		if err != nil {
//			// handle error
//		}
//
//		wd, err := os.Getwd()
//		if err != nil {
//			// handle error
//		}
//
//		buildPkg, err := build.ImportDir(".", wd, build.ImportComment)
//		if err != nil {
//			// handle error
//		}
//
//		// Create a documentation package from the build representation of our
//		// package.
//		pkg, err := lang.NewPackageFromBuild(buildPkg)
//		if err != nil {
//			// handle error
//		}
//
//		// Write the documentation out to console.
//		fmt.Println(out.Package(pkg))
//	}
package gomarkdoc
