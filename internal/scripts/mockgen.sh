#!/bin/bash

echo $file

fileBeforeDot=${file:0:${#file}-3}

firstLine=$(head -n 1 $file)

pkg=$(echo $firstLine | sed -e s/package\ //g)

mockgen -source $file -destination ${fileBeforeDot}_mock.go  -package $pkg
