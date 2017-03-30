#!/bin/bash

GOOS=windows go-bindata-assetfs -ignore='^(public/node_modules)' ./public/...

CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 go build -ldflags "-H windowsgui" -o Brook.386.exe .
