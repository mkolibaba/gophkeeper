package gen

import (
	"errors"
	"os"

	"github.com/mkolibaba/gophkeeper/shared/mage/tool"
	"github.com/uwu-tools/magex/mgx"
	"github.com/uwu-tools/magex/shx"
)

var (
	must = shx.CommandBuilder{StopOnError: true}
)

// Run mockery
func Mockery() {
	if _, err := os.Stat(".mockery.yml"); errors.Is(err, os.ErrNotExist) {
		mgx.Must(errors.New(".mockery.yml does not exist"))
	}

	tool.Install("mockery", "github.com/vektra/mockery/v3@latest")

	must.RunV("mockery")
}
