#!/bin/bash

for dir in */ ; do
  cd "$dir" || exit
  ./check.sh
  cd ..
done

