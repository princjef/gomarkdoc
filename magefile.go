//go:build mage
// +build mage

package main

import (
	"os"
	"path/filepath"

	"github.com/princjef/mageutil/bintool"
	"github.com/princjef/mageutil/shellcmd"
)

var linter = bintool.Must(bintool.New(
	"golangci-lint{{.BinExt}}",
	"1.45.2",
	"https://github.com/golangci/golangci-lint/releases/download/v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
))

func Lint() error {
	if err := linter.Ensure(); err != nil {
		return err
	}

	return linter.Command(`run`).Run()
}

func Generate() error {
	return shellcmd.Command(`go generate .`).Run()
}

func Build() error {
	return shellcmd.Command(`go build -o ./bin/gomarkdoc ./cmd/gomarkdoc`).Run()
}

func Doc() error {
	return shellcmd.RunAll(
		`go run ./cmd/gomarkdoc .`,
		`go run ./cmd/gomarkdoc --header "" ./lang/...`,
		`go run ./cmd/gomarkdoc --header "" ./format/...`,
		`go run ./cmd/gomarkdoc --header "" ./cmd/...`,
	)
}

func DocVerify() error {
	return shellcmd.RunAll(
		`go run ./cmd/gomarkdoc -c .`,
		`go run ./cmd/gomarkdoc -c --header "" ./lang/...`,
		`go run ./cmd/gomarkdoc -c --header "" ./format/...`,
		`go run ./cmd/gomarkdoc -c --header "" ./cmd/...`,
	)
}

func RegenerateTestDocs() error {
	dirs, err := os.ReadDir("./testData")
	if err != nil {
		return err
	}

	base, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		os.Chdir(filepath.Join(base, "./testData", dir.Name()))
		if err := shellcmd.Command(`go run ../../cmd/gomarkdoc -o "{{.Dir}}/README.md" ./...`).Run(); err != nil {
			return err
		}
	}

	return nil
}

func Test() error {
	return shellcmd.Command(`go test -count 1 -coverprofile=coverage.txt ./...`).Run()
}

func Coverage() error {
	return shellcmd.Command(`go tool cover -html=coverage.txt`).Run()
}
