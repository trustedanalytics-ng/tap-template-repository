#!/bin/bash -x
files=$(find  ./ -name *mock* |  grep -v vendor |  grep -v \.git)
for f in $files; do
  echo $f
  cat License.txt | cat - $f > temp && mv temp $f
done
