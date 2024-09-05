package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Application struct {
	command string
	args    []string
}

func (a Application) GetCommander() (Commander, error) {
	if a.command == "cat-file" {
		comm := BlobRead{args: a.args}
		return &comm, nil
	}
	if a.command == "init" {
		init := Init{}
		return &init, nil
	}
	return nil, errors.New("Command not found")
}

type Commander interface {
	Exec()
}

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}
	app := Application{command: os.Args[1], args: os.Args[2:]}
	action, err := app.GetCommander()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	action.Exec()
}
