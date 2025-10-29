package tool

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/uwu-tools/magex/shx"
	"os/exec"
)

var (
	must = shx.CommandBuilder{StopOnError: true}
)

func Install(tool string, link string) {
	if _, err := exec.LookPath(tool); err == nil {
		return
	}
	color.HiYellow(fmt.Sprintf("%s not found, installing...", tool))
	must.RunV("go", "install", link)
}
