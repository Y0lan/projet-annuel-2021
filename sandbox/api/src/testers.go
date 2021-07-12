package main

import (
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func testPython(code, test string) (output string, status bool) {

	status = SUCCESS
	dir := getTemporaryDirectoryForCodeExecution()
	defer cleanUpTemporaryFiles(dir)
	createTemporaryFile(code, "py", "main", "", "", dir)
	pythonTestFile := createTemporaryFile(test, "py", "", "", "", dir)

	out, err := exec.Command("python3", pythonTestFile.Name()).CombinedOutput()
	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, output)
	}
	return
}

// enableTestFileToAccessSubmittedCode replace EXERCISE_FILE_RANDOM by the name of the file so the test file can access the function
func enableTestFileToAccessSubmittedCode(test, moduleName string) string {
	var replacer = strings.NewReplacer("EXERCISE_FILE_RANDOM", strings.TrimSuffix(filepath.Base(moduleName), path.Ext(moduleName)))
	test = replacer.Replace(test)
	return test
}

func testRust(code, test string) (output string, status bool) {

	status = SUCCESS
	dir := getTemporaryDirectoryForCodeExecution()
	defer cleanUpTemporaryFiles(dir)
	rustCodeFile := createTemporaryFile(code, "rs", "", "code", "", dir)

	test = enableTestFileToAccessSubmittedCode(test, rustCodeFile.Name())

	rustTestFile := createTemporaryFile(test, "rs", "", "", "", dir)

	out, err := exec.Command("rustc", "--test", "--crate-name", "test", rustTestFile.Name()).CombinedOutput()
	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, output)
		return
	}

	out, err = exec.Command("./test").CombinedOutput()
	output = string(out)
	if err != nil {
		status = ShowAndSetError(err, output)
	}
	return
}

func testGo(code, test string) (output string, status bool) {

	status = SUCCESS
	dir := getTemporaryDirectoryForCodeExecution()
	defer cleanUpTemporaryFiles(dir)

	createTemporaryFile("module challenge\n\ngo 1.16\n", "mod", "go", "", "", dir)
	createTemporaryFile(code, "go", "", "", "", dir)
	createTemporaryFile(test, "go", "", "", "_test", dir)

	command := exec.Command("go", "test")
	command.Dir = dir
	out, err := command.CombinedOutput()

	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, output)
	}
	return
}
