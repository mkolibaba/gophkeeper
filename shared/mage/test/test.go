package test

import (
	"github.com/fatih/color"
	"github.com/uwu-tools/magex/shx"
	"strings"
)

var (
	must = shx.CommandBuilder{StopOnError: true}
)

// Running all tests with prettified output
func Run() error {
	color.HiGreen("Running tests...")

	output, err := shx.Output("go", "test", "./...")
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
