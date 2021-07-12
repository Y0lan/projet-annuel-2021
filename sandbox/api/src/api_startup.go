package main

import (
	"fmt"
	"log"
	"os/exec"
)

func displayIfEveryLangIsCorrectlyConfigured() {
	out, err := exec.Command("/home/runner/api/tests/test_all.sh").Output()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("ERROR: DOCKER IMAGE NOT CONFIGURED CORRECTLY. EXITING...")
	}
	fmt.Println(string(out))
}
