#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o brook .
GOOS=linux GOARCH=386 go build -o brook386 .
GOOS=darwin GOARCH=amd64 go build -o brookmacos .
GOOS=windows GOARCH=amd64 go build -o brook.exe .
GOOS=windows GOARCH=386 go build -o brook.386.exe .
