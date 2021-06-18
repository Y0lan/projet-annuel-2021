#!/bin/bash

echo "-------------------"
echo "testing PYTHON installation"
function run_program_and_test() {
  # run all python files
  for python_file in *.py; do
    if ! python3 "$python_file" > /dev/null 2> err.log; then
      cat err.log
      rm err.log
      echo "PYTHON KO"
      return 1
    fi
  done

  # remove logs files
  if [ -f "err.log" ]; then
    rm err.log
  fi
}

run_program_and_test
echo "PYTHON OK"
