package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os/exec"
	"runtime"
	"strings"
)

var (
	outPath        = "bin/gophkeeper-server"
	sqlcRepository = "github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0"
	sqlcConfig     = "sqlite/sqlc/sqlc.yml"
)

// Run gen
func Gen() {
	mg.Deps(GenSqlc)
}

// Generate sqlc
func GenSqlc() error {
	return sh.RunV("go", "run", sqlcRepository, "generate", "-f", sqlcConfig)
}

// Run server build
func Build() error {
	return sh.RunV("go", "build",
		"-ldflags", "-s -w",
		"-o", binaryPath(),
		"cmd/server/main.go")
}

// Running server
func Run() error {
	mg.Deps(Build)
	return sh.RunV("./" + binaryPath())
}

// Run tests
func Test() error {
	output, err := sh.Output("go", "test", "./...")
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "[no test files]") {
			continue
		}
		if strings.HasPrefix(line, "ok") {
			color.HiGreen(line)
		} else if strings.Contains(line, "FAIL") {
			color.HiRed(line)
		} else {
			color.New().Println(line)
		}
	}
	return err
}

// Run server in watch mode (requires watchexec)
func Watch() error {
	if _, err := exec.LookPath("watchexec"); err != nil {
		fmt.Println("please install watchexec (https://github.com/watchexec/watchexec) to use this target")
		return nil
	}

	return sh.RunV("watchexec",
		"-r",
		"-c",
		"-e", "go",
		"--poll", "1000",
		"--",
		"go", "run", "cmd/server/main.go")
}

func binaryPath() string {
	if runtime.GOOS == "windows" {
		return outPath + ".exe"
	}
	return outPath
}
