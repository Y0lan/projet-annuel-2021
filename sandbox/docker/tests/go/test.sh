#!/bin/bash

echo "-------------------"
echo "testing GO installation"

function show_errors() {
      cat err.log
      rm err.log
      echo "GO KO"
      exit 1
}

function run_program_and_test() {

    if ! go run main.go > /dev/null 2> err.log; then
      show_errors
    fi

    if ! go test > /dev/null 2> err.log; then
      show_errors
    fi

    # remove logs files
    if [ -f "err.log" ]; then
      rm err.log
    fi
}

run_program_and_test
echo "GO OK"
