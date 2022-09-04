#!/bin/sh

go build -o build/ main.go
build/main -file test.yaml
