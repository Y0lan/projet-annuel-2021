#!/bin/bash

test_dir="$HOME/api/tests"
for dir in "$test_dir"/*/ ; do
  dir=${dir%/}
  cd "$dir" || exit
  if ! ./test.sh; then
    exit 1
  fi
done

# all is healthy
exit 0

