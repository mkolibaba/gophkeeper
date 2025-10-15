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

// Run gen target for given module (running for all in none specified)
func Gen(module string) {
	forEachModule(func() {
		sh.RunV("mage", "gen")
	}, module)
}

// Run test for given module (running for all in none specified)
func Test(module string) {
	forEachModule(func() {
		sh.RunV("mage", "test")
	}, module)
}

// Run 'go mod tidy' for given module (running for all in none specified)
func Tidy(module string) {
	forEachModule(func() {
		sh.RunV("go", "mod", "tidy")
	}, module)
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
