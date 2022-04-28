package main

import (
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/LassiHeikkila/taskey/pkg/types"
)

type taskExecCallback func(name string, status int, output string)

func makeTask(task *types.Task, cb taskExecCallback) func() {
	switch task.Content.(type) {
	case *types.CmdTask:
		return makeCmdTask(task, cb)
	case *types.ScriptTask:
		return makeScriptTask(task, cb)
	default:
		return nil
	}
}

func makeCmdTask(task *types.Task, cb taskExecCallback) func() {
	return func() {
		cmdTask := task.Content.(*types.CmdTask)

		cmd := exec.Command(cmdTask.Program, cmdTask.Args...)
		var status int
		var output string
		err := execCmd(cmd, cmdTask.CombinedOutput, &status, &output)
		if err != nil {
			log.Println("failed to execute command task:", err)
			return
		}

		cb(task.Name, status, output)
	}
}

func makeScriptTask(task *types.Task, cb taskExecCallback) func() {
	return func() {
		scriptTask := task.Content.(*types.ScriptTask)
		f, _ := ioutil.TempFile("", "taskey-script") // TODO: check error
		_, _ = f.WriteString(scriptTask.Script)      // TODO: check error
		_ = f.Close()                                // TODO: check error

		cmd := exec.Command(scriptTask.Interpreter, f.Name())

		var status int
		var output string
		err := execCmd(cmd, scriptTask.CombinedOutput, &status, &output)
		if err != nil {
			log.Println("failed to execute script task:", err)
			return
		}

		cb(task.Name, status, output)
	}
}

func execCmd(cmd *exec.Cmd, combinedOutput bool, status *int, output *string) error {
	if status == nil {
		status = new(int)
	}
	if output == nil {
		output = new(string)
	}

	var b []byte
	var err error

	if combinedOutput {
		b, err = cmd.CombinedOutput()
	} else {
		b, err = cmd.Output()
	}

	*output = string(b)

	if err == nil {
		*status = 0
		return nil
	}

	switch err := err.(type) {
	case *exec.ExitError:
		*status = err.ExitCode()
		return nil
	default:
	}
	return err
}
