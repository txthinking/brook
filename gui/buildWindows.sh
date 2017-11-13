GOOS=windows go-bindata-assetfs -ignore='^(public/node_modules|public/dl)' ./public/...

cp iconwin.go.black iconwin.go
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows go build -ldflags "-H windowsgui" -o Brook.exe .

cp iconwin.go.white iconwin.go
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows go build -ldflags "-H windowsgui" -o Brook.white.exe .
