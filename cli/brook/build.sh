#!/bin/bash

if [ $# -ne 1 ]; then
    echo "./build.sh version"
    exit
fi

mkdir _

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build                    -ldflags="-w -s" -o _/brook_darwin_amd64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build                    -ldflags="-w -s" -o _/brook_darwin_arm64 .
CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build                     -ldflags="-w -s" -o _/brook_freebsd_386 .
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build                   -ldflags="-w -s" -o _/brook_freebsd_amd64 .
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build                     -ldflags="-w -s" -o _/brook_linux_amd64 .
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build                       -ldflags="-w -s" -o _/brook_linux_386 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build                     -ldflags="-w -s" -o _/brook_linux_arm64 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build               -ldflags="-w -s" -o _/brook_linux_arm7 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build               -ldflags="-w -s" -o _/brook_linux_arm6 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build               -ldflags="-w -s" -o _/brook_linux_arm5 .
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build                      -ldflags="-w -s" -o _/brook_linux_mips .
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build                    -ldflags="-w -s" -o _/brook_linux_mipsle .
CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build     -ldflags="-w -s" -o _/brook_linux_mips_softfloat .
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build   -ldflags="-w -s" -o _/brook_linux_mipsle_softfloat .
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build                    -ldflags="-w -s" -o _/brook_linux_mips64 .
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build                  -ldflags="-w -s" -o _/brook_linux_mips64le .
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 GOMIPS=softfloat go build   -ldflags="-w -s" -o _/brook_linux_mips64_softfloat .
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le GOMIPS=softfloat go build -ldflags="-w -s" -o _/brook_linux_mips64le_softfloat .
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64 go build                     -ldflags="-w -s" -o _/brook_linux_ppc64 .
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build                   -ldflags="-w -s" -o _/brook_linux_ppc64le .
CGO_ENABLED=0 GOOS=netbsd GOARCH=386 go build                      -ldflags="-w -s" -o _/brook_netbsd_386 .
CGO_ENABLED=0 GOOS=netbsd GOARCH=amd64 go build                    -ldflags="-w -s" -o _/brook_netbsd_amd64 .
CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go build                     -ldflags="-w -s" -o _/brook_openbsd_386 .
CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build                   -ldflags="-w -s" -o _/brook_openbsd_amd64 .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build                   -ldflags="-w -s" -o _/brook_windows_amd64.exe .
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build                     -ldflags="-w -s" -o _/brook_windows_386.exe .

cd ../../ipk

cp ../cli/brook/_/brook_linux_amd64 data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_amd64.ipk

cp ../cli/brook/_/brook_linux_386 data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_386.ipk

cp ../cli/brook/_/brook_linux_arm64 data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_arm64.ipk

cp ../cli/brook/_/brook_linux_arm7 data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_arm7.ipk

cp ../cli/brook/_/brook_linux_arm6 data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_arm6.ipk

cp ../cli/brook/_/brook_linux_arm5 data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_arm5.ipk

cp ../cli/brook/_/brook_linux_mips data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_mips.ipk

cp ../cli/brook/_/brook_linux_mipsle data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_mipsle.ipk

cp ../cli/brook/_/brook_linux_mips_softfloat data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_mips_softfloat.ipk

cp ../cli/brook/_/brook_linux_mipsle_softfloat data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_mipsle_softfloat.ipk

cp ../cli/brook/_/brook_linux_mips64 data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_mips64.ipk

cp ../cli/brook/_/brook_linux_mips64le data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_mips64le.ipk

cp ../cli/brook/_/brook_linux_mips64_softfloat data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_mips64_softfloat.ipk

cp ../cli/brook/_/brook_linux_mips64le_softfloat data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_mips64le_softfloat.ipk

cp ../cli/brook/_/brook_linux_ppc64 data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_ppc64.ipk

cp ../cli/brook/_/brook_linux_ppc64le data/usr/bin/brook
gsed -i 's/Architecture: .*/Architecture: all/' ./control/control
./build.sh
mv brook.ipk ../cli/brook/_/brook_linux_ppc64le.ipk

cd ../cli/brook/

nami release github.com/txthinking/brook $1 _

rm -rf _
