package main

import (
	"bytes"
	"fmt"
	"github.com/bitfield/script"
	"github.com/carolynvs/magex/shx"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os/exec"
	"runtime"
	"strings"
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

// Runs go test in verbose mode and prettifies the output
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
	mg.Deps(GenMock)
}

// Run tests with coverage
func TestCoverage() {
	installTool("go-test-coverage", "github.com/vladopajic/go-test-coverage/v2@latest")

	color.HiYellow("[testcoverage] Running tests...")
	must.RunV("go", "test", "./...", "-coverprofile=./cover.out", "-covermode=atomic", "-coverpkg=./...")

	color.HiYellow("[testcoverage] Creating html coverage file...")
	must.RunV("go", "tool", "cover", "-html", "cover.out", "-o", "cover.html")

	color.HiYellow("[testcoverage] Running go-test-coverage...")
	must.RunV("go-test-coverage", "--config=./.testcoverage.yml", "--badge-file-name=./coverage.svg")

	color.HiGreen("[testcoverage] Done")
}

func installTool(tool, link string) {
	if _, err := exec.LookPath(tool); err == nil {
		return
	}
	color.HiYellow(fmt.Sprintf("%s not found, installing...", tool))
	must.RunV("go", "install", link)
}

// Generate mocks
func GenMock() error {
	installTool("moq", "github.com/matryer/moq@latest")

	color.HiGreen("Generating mocks")
	return sh.RunV("go", "generate", "github.com/mkolibaba/gophkeeper/client")
}

func installMoq() error {
	if _, err := exec.LookPath("moq"); err == nil {
		return nil
	}

	fmt.Println("moq has not been found, installing")
	return sh.RunV("go", "install", "github.com/matryer/moq@latest")
}

func binaryPath() string {
	if runtime.GOOS == "windows" {
		return outPath + ".exe"
	}
	return outPath
}

func run(program string, args []string, env map[string]string) (string, error) {
	// Make string representation of command
	fullArgs := append([]string{program}, args...)
	cmdStr := strings.Join(fullArgs, " ")

	// Make string representation of environment
	envStrBuf := new(bytes.Buffer)
	for key, value := range env {
		fmt.Fprintf(envStrBuf, "%s=\"%s\", ", key, value)
	}
	envStr := string(bytes.TrimRight(envStrBuf.Bytes(), ", "))

	// Show info
	fmt.Println("Running '" + cmdStr + "'" + " with env " + envStr)

	// Run
	return sh.OutputWith(env, program, args...)
}
