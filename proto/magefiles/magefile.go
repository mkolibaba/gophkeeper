package main

import (
	"fmt"
	"github.com/bitfield/script"
	"os/exec"
)

var protoPath = "proto"

func Gen() error {
	// Убеждаемся, что protoc установлен.
	if _, err := exec.LookPath("protoc"); err != nil {
		return fmt.Errorf("please install protoc to use this target")
	}

	// Находим все .proto файлы
	files, err := script.ListFiles(protoPath + "/*.proto").Slice()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no .proto files found in %s", protoPath)
	}

	// Затем выполняем генерацию.
	// (просто указать последним аргументом proto/*.proto не получится, т.к. "* is a shell thing")
	_, err = script.Slice(files).
		ExecForEach("protoc -I=proto " +
			"--go_out=gen/go/gophkeeper " +
			"--go_opt=paths=source_relative " +
			"--go-grpc_out=gen/go/gophkeeper " +
			"--go-grpc_opt=paths=source_relative " +
			"--go_opt=default_api_level=API_OPAQUE {{.}}").
		Stdout()
	return err
}
