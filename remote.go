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

	// get called script
	if len(os.Args) < 2 {
		println("Error: No script name provided")
		println("Terminating.")
		return
	}

	var has_filename bool
	var filename string
	scriptName := ""

	skipNext := false

	for i, e := range os.Args[1:len(os.Args)] {
		if skipNext {
			skipNext = false
			continue
		}
		if e == "-h" || e == "--help" {
			println("Usage: remote [-h|--help] [-f|--file REMOTE FILE] command -- [passover arguments]")
			println("    -h | --help               : show help and exit")
			println("    -f | --file REMOTE FILE   : use REMOTE FILE as the script to run from; default to `remote.toml`")
			println("    command                   : the script name or command to run")
			println("    passover arguments        : arguments to pass over to the script, multiple -- can be used to chain passover arguments")
			return
		} else if (e == "-f" || e == "--file") && (i+1+1 < len(os.Args)) {
			has_filename = true
			filename = os.Args[i+1+1]
			skipNext = true
		} else if e == "--" {
			// passover args
			break
		} else {
			scriptName = os.Args[i+1]
		}
	}

	var pass_over_args []string
	var in_pass_over bool

	for _, e := range os.Args[1:len(os.Args)] {
		if e == "--" && !in_pass_over {
			in_pass_over = true
			continue
		}

		if in_pass_over {
			pass_over_args = append(pass_over_args, e)
		}
	}

	var scripts *toml.Tree

	if !has_filename {
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

		scripts = config.Get("scripts").(*toml.Tree)
	} else {
		_, err2 := os.Open(cwd + "/" + filename)

		if err2 != nil {
			println("Error: ", err2.Error())
			println("Cannot find or open `" + filename + "` in current working directory")
			println("Terminating.")
			return
		}

		config, err3 := toml.LoadFile(path.Join(cwd, filename))

		if err3 != nil {
			println("Error: ", err3.Error())
			println("Cannot parse `" + filename + "`")
			println("Terminating.")
			return
		}

		scripts = config.Get("scripts").(*toml.Tree)
	}

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
		execCommand(script.([]interface{}), pass_over_args)
	} else {
		// script is struct
		execCommandAndNext(script.(*toml.Tree), scripts, pass_over_args)
	}
}

func pop[T any](arr []T) (T, []T) {
	return arr[0], arr[1:]
}

func execCommand(command []interface{}, extra_args []string) {
	// convert to string array
	var args []string

	for _, arg := range command {
		args = append(args, arg.(string))
	}

	args = append(args, extra_args...)

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

func execCommandAndNext(command *toml.Tree, scripts *toml.Tree, extra_args []string) {
	// get command
	cmd := command.Get("command").([]interface{})

	// convert to string array
	var args []string

	for _, arg := range cmd {
		args = append(args, arg.(string))
	}

	// for this one, do not just copy, but quite literally pop it out
	// and then append the extra args
	for len(extra_args) > 0 {
		var arg string
		arg, extra_args = pop(extra_args)

		if arg == "--" {
			break
		}

		args = append(args, arg)
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
			execCommand(nextScript.([]interface{}), extra_args)
		} else {
			// next script is struct
			execCommandAndNext(nextScript.(*toml.Tree), scripts, extra_args)
		}
	}
}
