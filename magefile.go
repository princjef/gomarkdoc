//+build mage

package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const linterVersion = "1.21.0"

func Lint() error {
	if err := ensureLinter(); err != nil {
		return err
	}

	return pipedCmd("./bin/golangci-lint", "run")
}

func Doc() error {
	if err := pipedCmd("go", "run", "./cmd/gomarkdoc", "."); err != nil {
		return err
	}

	if err := pipedCmd("go", "run", "./cmd/gomarkdoc", "--header", "", "./lang/..."); err != nil {
		return err
	}

	if err := pipedCmd("go", "run", "./cmd/gomarkdoc", "--header", "", "./format/..."); err != nil {
		return err
	}

	if err := pipedCmd("go", "run", "./cmd/gomarkdoc", "--header", "", "./cmd/..."); err != nil {
		return err
	}

	return nil
}

func DocVerify() error {
	if err := pipedCmd("go", "run", "./cmd/gomarkdoc", "-c", "."); err != nil {
		return err
	}

	if err := pipedCmd("go", "run", "./cmd/gomarkdoc", "-c", "--header", "", "./lang/..."); err != nil {
		return err
	}

	if err := pipedCmd("go", "run", "./cmd/gomarkdoc", "-c", "--header", "", "./format/..."); err != nil {
		return err
	}

	if err := pipedCmd("go", "run", "./cmd/gomarkdoc", "-c", "--header", "", "./cmd/..."); err != nil {
		return err
	}

	return nil
}

func Test() error {
	return pipedCmd("go", "test", "-coverprofile=coverage.out", "./...")
}

func Coverage() error {
	return pipedCmd("go", "tool", "cover", "-html=coverage.out")
}

func ensureLinter() error {
	out, err := exec.Command("./bin/golangci-lint", "--version").Output()

	// If there was no error and we got the version we wanted, we can continue
	if err == nil && bytes.Contains(out, []byte(fmt.Sprintf(" %s ", linterVersion))) {
		return nil
	}

	// Install the linter if we don't have it
	downloadURL := fmt.Sprintf(
		"https://github.com/golangci/golangci-lint/releases/download/v%s/golangci-lint-%s-%s-%s.%s",
		linterVersion,
		linterVersion,
		runtime.GOOS,
		getArch(),
		getExt(),
	)

	res, err := http.Get(downloadURL)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("received unexpected response when downloading file: %d", res.StatusCode)
	}

	var name string
	var data []byte
	if getExt() == "tar.gz" {
		name, data, err = extractTarFile(res.Body, func(name string) bool {
			return strings.HasPrefix(name, "golangci-lint")
		})
		if err != nil {
			return err
		}
	} else {
		name, data, err = extractZipFile(res.Body, func(name string) bool {
			return strings.HasPrefix(name, "golangci-lint")
		})
		if err != nil {
			return err
		}
	}

	// Save the file
	if err := os.MkdirAll("./bin", 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("./bin/%s", name), data, 0755)
}

func extractTarFile(r io.Reader, fileTest func(name string) bool) (name string, data []byte, err error) {
	var buf bytes.Buffer

	zr, err := gzip.NewReader(r)
	if err != nil {
		return "", nil, err
	}

	tr := tar.NewReader(zr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return "", nil, errors.New("no executable found in archive")
		}

		if err != nil {
			return "", nil, err
		}

		name = filepath.Base(header.Name)

		if !fileTest(name) {
			continue
		}

		if _, err := io.Copy(&buf, tr); err != nil {
			return "", nil, err
		}

		break
	}

	return name, buf.Bytes(), nil
}

func extractZipFile(r io.Reader, fileTest func(name string) bool) (name string, data []byte, err error) {
	var rawBuf bytes.Buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&rawBuf, r); err != nil {
		return "", nil, err
	}

	zr, err := zip.NewReader(bytes.NewReader(rawBuf.Bytes()), int64(len(rawBuf.Bytes())))
	if err != nil {
		return "", nil, err
	}

	var found bool
	for _, f := range zr.File {
		name = filepath.Base(f.Name)

		if !fileTest(name) {
			continue
		}

		fc, err := f.Open()
		if err != nil {
			return "", nil, err
		}

		defer fc.Close()

		if _, err := io.Copy(&buf, fc); err != nil {
			return "", nil, err
		}

		found = true
		break
	}

	if !found {
		return "", nil, errors.New("no executable found in archive")
	}

	return name, buf.Bytes(), nil
}

func getExt() string {
	if runtime.GOOS == "windows" {
		return "zip"
	} else {
		return "tar.gz"
	}
}

func getArch() string {
	switch runtime.GOARCH {
	case "arm":
		return "armv7"
	default:
		return runtime.GOARCH
	}
}

func pipedCmd(name string, args ...string) error {
	fmt.Printf("%s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
