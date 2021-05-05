#!/bin/bash
mkdir -p build
GOOS=linux GOARCH=arm GOARM=5 go build -o build/$(basename $PWD) .
scp build/$(basename $PWD) taxi105:/home/root/software_go/