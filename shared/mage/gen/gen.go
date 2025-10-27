package gen

import (
	"errors"
	"github.com/mkolibaba/gophkeeper/shared/mage/tool"
	"github.com/uwu-tools/magex/mgx"
	"github.com/uwu-tools/magex/shx"
	"os"
)

var (
	must = shx.CommandBuilder{StopOnError: true}
)

// Run mockery
func Mockery() {
	if _, err := os.Stat(".mockery.yml"); errors.Is(err, os.ErrNotExist) {
		mgx.Must(errors.New(".mockery.yml does not exist"))
	}

	tool.Install("mockery", "github.com/vektra/mockery/v3@v3.5.5")

	must.RunV("mockery")
}
