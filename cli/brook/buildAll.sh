#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o brook .
GOOS=linux GOARCH=386 go build -o brook_linux_386 .
GOOS=linux GOARCH=arm64 go build -o brook_linux_arm64 .
GOOS=linux GOARCH=arm GOARM=7 go build -o brook_linux_arm_7 .
GOOS=linux GOARCH=arm GOARM=6 go build -o brook_linux_arm_6 .
GOOS=linux GOARCH=arm GOARM=5 go build -o brook_linux_arm_5 .
GOOS=darwin GOARCH=amd64 go build -o brook_macos_amd64 .
GOOS=windows GOARCH=amd64 go build -o brook_windows_amd64.exe .
GOOS=windows GOARCH=386 go build -o brook_windows_386.exe .
