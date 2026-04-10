package targets

import (
	"fmt"
	"os/exec"

	"github.com/bitfield/script"
	"github.com/magefile/mage/sh"
	"github.com/uwu-tools/magex/mgx"
	"github.com/uwu-tools/magex/pkg"
)

var protoPath = "proto"

// Запускает генерацию proto файлов
func Gen() error {
	// Убеждаемся, что protoc, protoc-gen-go и protoc-gen-go-grpc установлены.
	if _, err := exec.LookPath("protoc"); err != nil {
		return fmt.Errorf("please install protoc to use this target")
	}

	mgx.Must(pkg.EnsurePackageWith(pkg.EnsurePackageOptions{
		Name:           "google.golang.org/protobuf/cmd/protoc-gen-go",
		VersionCommand: "--version",
	}))
	mgx.Must(pkg.EnsurePackageWith(pkg.EnsurePackageOptions{
		Name:           "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		VersionCommand: "--version",
	}))

	// Находим все .proto файлы
	files, err := script.ListFiles(protoPath + "/*.proto").Slice()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no .proto files found in %s", protoPath)
	}

	// Затем выполняем генерацию.
	args := []string{
		"--proto_path=proto",
		"--go_out=gen/go/gophkeeperv1",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=gen/go/gophkeeperv1",
		"--go-grpc_opt=paths=source_relative",
		"--go_opt=default_api_level=API_OPAQUE",
	}
	args = append(args, files...)

	return sh.RunV("protoc", args...)
}
