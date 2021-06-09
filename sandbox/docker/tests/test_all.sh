#!/bin/bash

for dir in */ ; do
  cd "$dir" || exit
  ./test.sh
  cd ..
done

