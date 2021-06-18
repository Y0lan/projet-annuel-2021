#!/bin/bash

test_dir="$HOME/shared/tests"
for dir in "$test_dir"/*/ ; do
  dir=${dir%/}
  cd "$dir" || exit
  ./test.sh
done

