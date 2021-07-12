package main

import "os/exec"

func compilePython(code string) (output string, status bool) {

	status = SUCCESS
	dir := getTemporaryDirectoryForCodeExecution()
	defer cleanUpTemporaryFiles(dir)

	pythonFile := createTemporaryFile(code, "py", "", "", "", dir)

	out, err := exec.Command("python3", pythonFile.Name()).CombinedOutput()
	output = string(out)
	if err != nil {
		status = ShowAndSetError(err, output)
	}
	return
}

func compileRust(code string) (output string, status bool) {

	status = SUCCESS
	dir := getTemporaryDirectoryForCodeExecution()
	defer cleanUpTemporaryFiles(dir)
	rustFile := createTemporaryFile(code, "rs", "", "", "", dir)

	out, err := exec.Command("rustc", "--crate-name", "binary", rustFile.Name()).CombinedOutput()
	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, output)
		return
	}

	out, err = exec.Command("./binary").CombinedOutput()
	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, output)
	}

	return
}

func compileGo(code string) (output string, status bool) {

	status = SUCCESS
	dir := getTemporaryDirectoryForCodeExecution()
	defer cleanUpTemporaryFiles(dir)
	goFile := createTemporaryFile(code, "go", "", "", "", dir)

	out, err := exec.Command("go", "run", goFile.Name()).CombinedOutput()
	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, output)
	}
	return
}
