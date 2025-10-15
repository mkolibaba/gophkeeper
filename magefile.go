//go:build mage
// +build mage

package main

import (
	"fmt"
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

// Run gen for all modules
func Gen() {
	forEachModule(func() {
		sh.RunV("mage", "gen")
	})
}

// Run tests for all modules
func Test() {
	forEachModule(func() {
		sh.RunV("go", "test", "./...")
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

func forEachModule(executor func()) {
	for _, module := range modules {
		fmt.Println("󰠱 " + module)
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
