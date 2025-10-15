package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	"os/exec"
)

var modules = []string{
	"client",
	"proto",
	"server",
}

// Run gen target for all modules
func Gen() {
	forEachModule(func() {
		sh.RunV("mage", "gen")
	})
}

// Run test for all modules
func Test() {
	forEachModule(func() {
		sh.RunV("mage", "test")
	})
}

// Run 'go mod tidy' for all modules
func Tidy() {
	forEachModule(func() {
		sh.RunV("go", "mod", "tidy")
	})
}

// Install mage
func EnsureMage() error {
	fmt.Println("Installing Mage")
	return exec.Command("go", "install", "github.com/magefile/mage@latest").Run()
}

func forEachModule(executor func(), ms ...string) {
	if len(ms) == 0 {
		ms = modules
	}

	for _, module := range ms {
		color.HiMagenta("󰠱 " + module)
		os.Chdir(module)
		executor()
		os.Chdir("..")
	}
}

func must(err error) {
	if err != nil {
		panic(mg.Fatal(1, err))
	}
}
