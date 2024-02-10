#!/bin/bash

echo '# Brook' > ../readme.md
echo '<!--SIDEBAR-->' >> ../readme.md
echo '<!--G-R3M673HK5V-->' >> ../readme.md
echo 'A cross-platform programmable network tool.' >> ../readme.md
echo '' >> ../readme.md
echo '# Sponsor' >> ../readme.md
echo '**❤️  [Shiliew - China Optimized Network App](https://www.txthinking.com/shiliew.html)**' >> ../readme.md

cat getting-started.md >> ../readme.md
cat gui.md >> ../readme.md
cat resources.md >> ../readme.md

echo '# CLI Documentation' >> ../readme.md
cd ../cli/brook
go build
mv brook ~/.nami/bin/
cd ../../docs
jb '$1`brook mdpage`.split("\n").filter(v=>!v.startsWith("[")).join("\n").replace("```\n```", "```\nbrook [全局参数] 子命令 [子命令参数]\n```").split("\n").forEach(v=> echo(v.startsWith("**") && !v.startsWith("**Usage") ? "- "+v : v))' >> ../readme.md

cat example.md >> ../readme.md
cat diagram.md >> ../readme.md

markdown ../readme.md ./index.html

echo '# Brook' > _.md
echo 'A cross-platform programmable network tool' >> _.md
echo '' >> _.md
echo '# Sponsor' >> _.md
echo '**❤️  [Shiliew - China Optimized Network App](https://www.txthinking.com/shiliew.html)**' >> _.md
mdtoc ../readme.md >> _.md
cat ../readme.md >> _.md
mv _.md ../readme.md

