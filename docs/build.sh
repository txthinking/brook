#!/bin/bash

echo '# Brook' > ../readme.md
echo '<!--SIDEBAR-->' >> ../readme.md
echo '<!--G-R3M673HK5V-->' >> ../readme.md
echo 'A cross-platform programmable network tool. 一个跨平台可编程网络工具' >> ../readme.md
echo '' >> ../readme.md
echo '# Sponsor' >> ../readme.md
echo '**❤️  [Shiliew - China Optimized VPN](https://www.txthinking.com/shiliew.html)**' >> ../readme.md

cat getting-started.md >> ../readme.md

cat install-cli.md >> ../readme.md
cat daemon.md >> ../readme.md
cat auto-start.md >> ../readme.md
cat one-click-script.md >> ../readme.md

echo '# CLI Documentation 命令行文档' >> ../readme.md
cd ../cli/brook
go build
mv brook ~/.nami/bin/
cd ../../docs
brook mdpage >> ../readme.md
cat gui.md >> ../readme.md
cat gui-zh.md >> ../readme.md
cat diagram.md >> ../readme.md

echo '# Protocol' >> ../readme.md
echo 'https://github.com/txthinking/brook/tree/master/protocol' >> ../readme.md
echo '# Blog' >> ../readme.md
echo 'https://www.txthinking.com/talks/' >> ../readme.md
echo '# YouTube' >> ../readme.md
echo 'https://www.youtube.com/txthinking' >> ../readme.md
echo '# Telegram' >> ../readme.md
echo 'https://t.me/s/txthinking_news' >> ../readme.md
echo '# brook-mamanger' >> ../readme.md
echo 'https://github.com/txthinking/brook-manager' >> ../readme.md
echo '# nico' >> ../readme.md
echo 'https://github.com/txthinking/nico' >> ../readme.md
echo '# Brook Deploy' >> ../readme.md
echo 'https://www.txthinking.com/deploy.html' >> ../readme.md
echo '# Pastebin' >> ../readme.md
echo 'https://paste.brook.app' >> ../readme.md
echo '# 独立脚本例子 | Standalone Script Example' >> ../readme.md
echo 'https://github.com/txthinking/bypass' >> ../readme.md
echo '# 脚本生成器 | Brook Script Builder' >> ../readme.md
echo 'https://modules.brook.app' >> ../readme.md

markdown ../readme.md ./index.html

echo '# Brook' > _.md
echo 'A cross-platform programmable network tool. 一个跨平台可编程网络工具' >> _.md
echo '' >> _.md
echo '# Sponsor' >> _.md
echo '**❤️  [Shiliew - China Optimized VPN](https://www.txthinking.com/shiliew.html)**' >> _.md
mdtoc ../readme.md >> _.md
cat ../readme.md >> _.md
mv _.md ../readme.md

