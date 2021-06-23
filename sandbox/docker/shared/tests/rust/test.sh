#!/bin/bash

echo "-------------------"
echo "testing RUST installation"

function show_errors() {
      cat err.log
      rm err.log
      echo "RUST KO"
      exit 1
}

function run_program_and_test() {
    # compile
    if ! cargo run > /dev/null 2> err.log ; then
      show_errors
    fi

    if ! cargo test > /dev/null 2> err.log ; then
      show_errors
    fi

    # remove logs files and binary file
    if [ -f "err.log" ]; then
      rm err.log
    fi
    rm -rf target
    rm Cargo.lock
    return 0
}

run_program_and_test
echo "RUST OK"
exit 0
