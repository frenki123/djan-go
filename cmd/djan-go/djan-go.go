package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed templates/main.tmpl
var mainTmpl []byte

func parseArgs(args []string) error {
	if len(args) == 1 {
		return errors.New("missing command arguments")
	}
	cmd := args[1]
	switch cmd {
	case "run":
		fmt.Println("Run command")
		return nil
	case "build":
		err := buildTemplate(args)
		return err
	}

	return errors.New("non existing command argument")
}

func buildTemplate(args []string) error {
	const filePerm = 0644
	const dirPerm = 0755
	const dirCreateErr = "error during creating directory '%s': %w"
	const fileWriteErr = "error during writing file '%s': %w"
	if len(args) < 3 {
		return errors.New("missing project name")
	}
	app := args[2]
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	projPath := filepath.Join(wd, app)
	if err := os.Mkdir(projPath, dirPerm); err != nil {
		return fmt.Errorf(dirCreateErr, projPath, err)
	}
	cmd := exec.Command("go", "mod", "init", "-C", projPath, fmt.Sprintf("example/%s", app))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go module not initilized: %w", err)
	}
	directories := []string{"cmd", "web", "internal"}
	for _, dir := range directories {
		curDir := filepath.Join(projPath, dir)
		if err := os.Mkdir(curDir, dirPerm); err != nil {
			return fmt.Errorf(dirCreateErr, curDir, err)
		}
	}
	cmdAppPath := filepath.Join(projPath, "cmd", app)
	if err := os.Mkdir(cmdAppPath, dirPerm); err != nil {
		return fmt.Errorf(dirCreateErr, cmdAppPath, err)
	}
	mainFile := filepath.Join(cmdAppPath, "main.go")
	if err := os.WriteFile(mainFile, mainTmpl, filePerm); err != nil {
		return fmt.Errorf(fileWriteErr, mainFile, err)
	}
	return nil
}

func main() {
	args := os.Args
	if err := parseArgs(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}
