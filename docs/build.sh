#!/bin/bash

echo '# Brook' > ../readme.md
echo '<!--SIDEBAR-->' >> ../readme.md
echo '<!--G-R3M673HK5V-->' >> ../readme.md
echo 'A cross-platform network tool' >> ../readme.md

cat getting-started.md >> ../readme.md

cat install-cli.md >> ../readme.md
cat daemon.md >> ../readme.md
cat auto-start.md >> ../readme.md
cat one-click-script.md >> ../readme.md

cat install-gui.md >> ../readme.md

echo '# CLI Documentation 命令行文档' >> ../readme.md
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
echo '# Pastebin' >> ../readme.md
echo 'https://ooo.soso.ooo' >> ../readme.md

markdown ../readme.md ./index.html

echo '# Brook' > _.md
echo 'A cross-platform network tool' >> _.md
mdtoc ../readme.md >> _.md
cat ../readme.md >> _.md
mv _.md ../readme.md

