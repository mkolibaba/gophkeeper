package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"os/exec"
	"runtime"
	"strings"
)

var (
	outPath        = "bin/gophkeeper-server"
	sqlcConfigPath = "sqlite/sqlc/sqlc.yml"
)

// Build server binary
func Build() error {
	mg.Deps(Gen)
	color.HiYellow("[build] Building server binary")
	if err := sh.RunV("go", "build",
		"-ldflags", "-s -w",
		"-o", binaryPath(),
		"cmd/server/main.go"); err != nil {
		return err
	}
	color.HiGreen("[build] Done")
	return nil
}

// Run server
func Run() error {
	mg.Deps(Build)
	color.HiGreen("Starting server...")
	return sh.RunV("./" + binaryPath())
}

// Runs go test in verbose mode and prettifies the output
func Test() error {
	color.HiGreen("Running tests...")
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
		"-e", "go,sql",
		"--wrap-process", "session",
		"--poll", "1000",
		"--",
		"mage", "run")
}

// Run all code generators
func Gen() {
	color.HiYellow("[gen] Generating sources")
	mg.Deps(GenSqlc, GenGoverter, GenOpaqueMapper)
	color.HiGreen("[gen] Done")
}

// Generate sqlc code from queries
func GenSqlc() error {
	var needsRefresh bool
	dsts := []string{"db.go", "models.go", "query.sql.go"}
	for _, dst := range dsts {
		var err error
		needsRefresh, err = target.Path(
			fmt.Sprintf("sqlite/sqlc/gen/%s", dst),
			"sqlite/sqlc/query.sql",
			sqlcConfigPath,
		)
		if err != nil {
			color.HiRed("[sqlc] Error: %s", err.Error())
			return nil
		}
		if needsRefresh {
			break
		}
	}

	if !needsRefresh {
		color.HiGreen("[sqlc] Generated files are up to date, skipping")
		return nil
	}

	mg.Deps(installSqlc)

	color.HiYellow("[sqlc] Generating code...")
	if err := sh.RunV("sqlc", "generate", "-f", sqlcConfigPath); err != nil {
		return err
	}

	color.HiGreen("[sqlc] Done")
	return nil
}

func GenGoverter() error {
	needsRefresh, err := target.Path(
		"sqlite/converter/gen/converter.go",
		"sqlite/converter/converter.go",
		"sqlite/sqlc/gen/models.go",
		"data.go",
	)
	if err != nil {
		return err
	}
	if !needsRefresh {
		color.HiGreen("[goverter] Generated files are up to date, skipping")
		return nil
	}

	mg.Deps(installGoverter)

	color.HiYellow("[goverter] Generating code...")
	if err := sh.RunV("goverter", "gen", "./sqlite/converter"); err != nil {
		return err
	}
	color.HiGreen("[goverter] Done")
	return nil
}

func GenOpaqueMapper() error {
	var needsRefresh bool
	dsts := []string{"binary", "card", "login", "note"}
	for _, dst := range dsts {
		var err error
		needsRefresh, err = target.Path(
			fmt.Sprintf("grpc/gen/%s_mapping.go", dst),
			"data.go",
		)
		if err != nil {
			color.HiRed("[opaquemapper] Error: %s", err.Error())
			return nil
		}
		if needsRefresh {
			break
		}
	}

	if !needsRefresh {
		color.HiGreen("[opaquemapper] Generated files are up to date, skipping")
		return nil
	}

	color.HiYellow("[opaquemapper] Generating code...")
	if err := sh.RunV("go", "generate", "data.go"); err != nil {
		return err
	}
	color.HiGreen("[opaquemapper] Done")
	return nil
}

func installGoverter() error {
	return installTool("goverter", "github.com/jmattheis/goverter/cmd/goverter@v1.9.1")
}

func installSqlc() error {
	return installTool("sqlc", "github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0")
}

func installTool(tool, link string) error {
	if _, err := exec.LookPath(tool); err == nil {
		return nil
	}
	color.HiYellow(fmt.Sprintf("%s not found, installing...", tool))
	return sh.RunV("go", "install", link)
}

func binaryPath() string {
	if runtime.GOOS == "windows" {
		return outPath + ".exe"
	}
	return outPath
}
