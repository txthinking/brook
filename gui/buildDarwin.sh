GOOS=darwin go-bindata-assetfs -ignore='^(public/node_modules|public/dl)' ./public/...

cp iconunix.go.black iconunix.go
CGO_ENABLED=1 GOOS=darwin go build -ldflags -s -o brook .
mv brook Brook.app/Contents/MacOS/
7z a Brook.app.zip Brook.app

cp iconunix.go.white iconunix.go
CGO_ENABLED=1 GOOS=darwin go build -ldflags -s -o brook .
mv brook Brook.app/Contents/MacOS/
7z a Brook.app.white.zip Brook.app
