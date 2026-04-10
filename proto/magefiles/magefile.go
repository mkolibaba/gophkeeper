package main

import (
	"fmt"
	"os/exec"

	"github.com/bitfield/script"
	"github.com/magefile/mage/sh"

	"github.com/mkolibaba/gophkeeper/shared/mage/tool"
)

var protoPath = "proto"

func Gen() error {
	// Убеждаемся, что protoc и protoc-gen-go установлены.
	if _, err := exec.LookPath("protoc"); err != nil {
		return fmt.Errorf("please install protoc to use this target")
	}

	tool.Install("protoc-gen-go", "google.golang.org/protobuf/cmd/protoc-gen-go@latest")
	tool.Install("protoc-gen-go-grpc", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")

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
