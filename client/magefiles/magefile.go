package main

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mkolibaba/gophkeeper/shared/mage/tool"
	"github.com/uwu-tools/magex/shx"
	"os/exec"
	"runtime"
	//mage:import test
	_ "github.com/mkolibaba/gophkeeper/shared/mage/test"
)

var (
	outPath = "bin/gophkeeper-client"
	logOut  = "bin/client.log"
	spewOut = "bin/spew.log"

	must = shx.CommandBuilder{StopOnError: true}
)

// Run client build
func Build() error {
	return sh.RunV("go", "build",
		"-ldflags", "-s -w",
		"-o", binaryPath(),
		"cmd/client/main.go")
}

// Running client
func Run() error {
	mg.Deps(Build)
	return sh.RunV("./" + binaryPath())
}

// Run client in watch mode (requires watchexec)
func Watch() error {
	if _, err := exec.LookPath("watchexec"); err != nil {
		fmt.Println("please install watchexec (https://github.com/watchexec/watchexec) to use this target")
		return nil
	}

	return sh.RunV("watchexec",
		"-r",
		"-c",
		"-e", "go",
		"--wrap-process", "session",
		"--poll", "1000",
		"--",
		"go", "run", "cmd/client/main.go")
}

// Show client log
func Log() error {
	_, err := script.File(logOut).Last(1000).Stdout()
	return err
}

// Show spew log
func Spew() error {
	_, err := script.File(spewOut).Last(1000).Stdout()
	return err
}

// Run gen
func Gen() {
	mg.Deps(GenMockery)
}

func GenMockery() {
	tool.Install("mockery", "github.com/vektra/mockery/v3@v3.5.5")
	must.RunV("mockery")
}

func binaryPath() string {
	if runtime.GOOS == "windows" {
		return outPath + ".exe"
	}
	return outPath
}
