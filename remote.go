package main

import (
	"os"
	"os/exec"
	"path"

	"github.com/pelletier/go-toml"
)

func main() {
	cwd, err := os.Getwd()

	if err != nil {
		println("Error: ", err.Error())
		println("Cannot get current working directory")
		println("Terminating.")
		return
	}

	// check if the current working directory contains a file named "remote.toml"

	_, err2 := os.Open(cwd + "/remote.toml")

	if err2 != nil {
		println("Error: ", err2.Error())
		println("Cannot find or open `remote.toml` in current working directory")
		println("Terminating.")
		return
	}

	// parse the file
	config, err3 := toml.LoadFile(path.Join(cwd, "remote.toml"))

	if err3 != nil {
		println("Error: ", err3.Error())
		println("Cannot parse `remote.toml`")
		println("Terminating.")
		return
	}

	scripts := config.Get("scripts").(*toml.Tree)

	// get called script
	if len(os.Args) < 2 {
		println("Error: No script name provided")
		println("Terminating.")
		return
	}

	scriptName := os.Args[1]

	// check if script exists
	if !scripts.Has(scriptName) {
		println("Error: Script `" + scriptName + "` does not exist")
		println("Terminating.")
		return
	}

	// get script as array or struct
	script := scripts.Get(scriptName)

	// check if script is array
	if _, ok := script.([]interface{}); ok {
		execCommand(script.([]interface{}))
	} else {
		// script is struct
		execCommandAndNext(script.(*toml.Tree), scripts)
	}
}

func execCommand(command []interface{}) {
	// convert to string array
	var args []string

	for _, arg := range command {
		args = append(args, arg.(string))
	}

	// execute command
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		println("Error: ", err.Error())
		println("Terminating.")
		return
	}
}

func execCommandAndNext(command *toml.Tree, scripts *toml.Tree) {
	// get command
	cmd := command.Get("command").([]interface{})

	// convert to string array
	var args []string

	for _, arg := range cmd {
		args = append(args, arg.(string))
	}

	// check if next exists
	next := command.Get("next")

	if next != nil {
		// get next as string
		next := command.Get("next").(string)

		// check if next is empty
		if next == "" {
			println("Error: `next` is empty")
			println("Terminating.")
			return
		}

		// check if next is a script that exists
		if !scripts.Has(next) {
			println("Error: Script `" + next + "` does not exist")
			println("Terminating.")
			return
		}
	}

	// execute command
	cmdHandle := exec.Command(args[0], args[1:]...)
	cmdHandle.Stdout = os.Stdout
	cmdHandle.Stderr = os.Stderr
	err := cmdHandle.Run()

	if err != nil {
		println("Error: ", err.Error())
		println("Terminating.")
		return
	}

	if next != nil {
		next := next.(string)
		// get next script
		nextScript := scripts.Get(next)

		// check if next script is array
		if _, ok := nextScript.([]interface{}); ok {
			execCommand(nextScript.([]interface{}))
		} else {
			// next script is struct
			execCommandAndNext(nextScript.(*toml.Tree), scripts)
		}
	}
}
