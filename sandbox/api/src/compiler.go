package main

import (
	"os/exec"
)

func compilePython(code string, command chan Command) {
	result := Command{
		output: "",
		status: SUCCESS,
	}
	dir := getTemporaryDirectoryForCodeExecution()
	defer cleanUpTemporaryFiles(dir)

	pythonFile := createTemporaryFile(code, "py", "", "", "", dir)

	out, err := exec.Command("python3", pythonFile.Name()).CombinedOutput()
	if err != nil {
		result.status = ShowAndSetError(err, string(out))
	}
	result.output = string(out)
	command <- result
}

func compileRust(code string, command chan Command) {
	result := Command{
		output: "",
		status: SUCCESS,
	}
	dir := getTemporaryDirectoryForCodeExecution()
	defer cleanUpTemporaryFiles(dir)
	rustFile := createTemporaryFile(code, "rs", "", "", "", dir)

	out, err := exec.Command("rustc", "--crate-name", "binary", rustFile.Name()).CombinedOutput()

	if err != nil {
		result.status = ShowAndSetError(err, string(out))
	}

	out, err = exec.Command("./binary").CombinedOutput()

	if err != nil {
		result.status = ShowAndSetError(err, string(out))
	}
	result.output = string(out)
	command <- result
}

func compileGo(code string, command chan Command) {

	result := Command{
		output: "",
		status: SUCCESS,
	}

	dir := getTemporaryDirectoryForCodeExecution()
	defer cleanUpTemporaryFiles(dir)
	goFile := createTemporaryFile(code, "go", "", "", "", dir)

	out, err := exec.Command("go", "run", goFile.Name()).CombinedOutput()

	if err != nil {
		result.status = ShowAndSetError(err, string(out))
	}
	result.output = string(out)
	command <- result
}
