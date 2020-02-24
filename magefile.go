//+build mage

package main

import (
	"github.com/princjef/mageutil/bintool"
	"github.com/princjef/mageutil/shellcmd"
)

var linter = bintool.Must(bintool.New(
	"golangci-lint{{.BinExt}}",
	"1.23.6",
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

func Test() error {
	return shellcmd.Command(`go test -coverprofile=coverage.out ./...`).Run()
}

func Coverage() error {
	return shellcmd.Command(`go tool cover -html=coverage.out`).Run()
}
