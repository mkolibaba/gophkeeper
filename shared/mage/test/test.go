package test

import (
	"github.com/fatih/color"
	"github.com/mkolibaba/gophkeeper/shared/mage/tool"
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

// Run test coverage
func Coverage() {
	tool.Install("go-test-coverage", "github.com/vladopajic/go-test-coverage/v2@latest")

	color.HiYellow("[testcoverage] Running tests...")
	must.RunV("go", "test", "./...", "-coverprofile=./cover.out", "-covermode=atomic", "-coverpkg=./...")

	color.HiYellow("[testcoverage] Creating html coverage file...")
	must.RunV("go", "tool", "cover", "-html", "cover.out", "-o", "cover.html")

	color.HiYellow("[testcoverage] Running go-test-coverage...")
	must.RunV("go-test-coverage", "--config=./.testcoverage.yml", "--badge-file-name=./coverage.svg")

	color.HiGreen("[testcoverage] Done")
}
