package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/main.tmpl
var mainTmpl string

//go:embed templates/urls.tmpl
var urlTmpl string

//go:embed templates/handlers.tmpl
var handlerTmpl string

func parseArgs(args []string) error {
	if len(args) == 1 {
		return errors.New("missing command arguments")
	}
	cmd := args[1]
	switch cmd {
	case "run":
		fmt.Println("Run command")
		return nil
	case "createproject":
		err := createProject(args)
		return err
	}

	return errors.New("non existing command argument")
}

func getAppNames(arg string) (string, string) {
	names := strings.Split(arg, "/")
	l := len(names)
	appName := names[l-1]
	if l == 1 {
		hostname := "github.com"
		username := "user"
		user, err := user.Current()
		if err == nil {
			username = user.Username
		}
		module := fmt.Sprintf("%s/%s/%s", hostname, username, appName)
		return appName, module
	}
	return appName, arg
}

func createProject(args []string) error {
	const dirPerm = 0755
	const dirCreateErr = "error during creating directory '%s': %w"
	const fileWriteErr = "error during writing file '%s': %w"
	if len(args) < 3 {
		return errors.New("missing project name")
	}
	app, module := getAppNames(args[2])
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	projPath := filepath.Join(wd, app)
	if err := os.Mkdir(projPath, dirPerm); err != nil {
		return fmt.Errorf(dirCreateErr, projPath, err)
	}
	dirs := map[string]string{"cmd": "", "handlers": "", "models": "", "router": ""}
	for dir := range dirs {
		curDir := filepath.Join(projPath, dir)
		if err := os.Mkdir(curDir, dirPerm); err != nil {
			return fmt.Errorf(dirCreateErr, curDir, err)
		}
		dirs[dir] = curDir
	}
	mainFile := filepath.Join(dirs["cmd"], "main.go")
	mainData := map[string]string{
		"Host":   "localhost",
		"Port":   "9999",
		"Module": module,
	}
	if err := writeTemplate(mainFile, mainTmpl, mainData); err != nil {
		return fmt.Errorf(fileWriteErr, mainFile, err)
	}
	urlFile := filepath.Join(dirs["router"], "urls.go")
	if err := writeTemplate(urlFile, urlTmpl, mainData); err != nil {
		return fmt.Errorf(fileWriteErr, urlFile, err)
	}
	handlerFile := filepath.Join(dirs["handlers"], "index.go")
	if err := writeTemplate(handlerFile, handlerTmpl, mainData); err != nil {
		return fmt.Errorf(fileWriteErr, urlFile, err)
	}
	initCmd := exec.Command("go", "mod", "init", "-C", projPath, module)
	getCmd := exec.Command("go", "get", "-C", projPath, "github.com/frenki123/djan-go")
	tidyCmd := exec.Command("go", "mod", "tidy", "-C", projPath)

	cmds := map[string]exec.Cmd{"go mod init": *initCmd, "go get djan-go": *getCmd, "go mod tidy": *tidyCmd}
	for name, cmd := range cmds {
		fmt.Println(name)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command '%s' failed: %w", name, err)
		}
	}
	return nil
}

func writeTemplate(filepath string, tmplString string, data any) error {
	t := template.New("main")
	t, err := t.Parse(string(tmplString))
	if err != nil {
		return err
	}
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	if err := t.Execute(file, data); err != nil {
		return err
	}
	return nil
}

func main() {
	args := os.Args
	if err := parseArgs(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}
