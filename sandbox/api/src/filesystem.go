package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// createTemporaryFile takes:
// content to write into the file
// extension of the file
// fileName (optional)  so that the file is not created with a random name
// prefix (optional) to put BEFORE the random file name
// suffix (optional) to put AFTER the random file name and BEFORE the extension
func createTemporaryFile(content, extension, fileName, prefix, suffix, dir string) os.File {
	binaryContent := []byte(content)
	var tmpFile *os.File
	var err error

	if fileName == "" {
		tmpFile, err = ioutil.TempFile(dir, prefix+"*"+suffix+"."+extension)
	} else {
		fileName = filepath.Join(dir, fileName+"."+extension)
		tmpFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	}
	if err != nil {
		log.Fatal(err)
	}

	defer closeFile(*tmpFile)
	if _, err := tmpFile.Write(binaryContent); err != nil {
		log.Fatal(err)
	}
	return *tmpFile
}

func closeFile(file os.File) {
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

func getTemporaryDirectoryForCodeExecution() string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err.Error())
	}
	return dir
}

func cleanUpTemporaryFiles(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatal("Could not remove directory after completion")
	}
}
