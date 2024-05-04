#!/bin/bash

echo '# Brook' > ../readme.md
echo '<!--SIDEBAR-->' >> ../readme.md
echo '<!--G-R3M673HK5V-->' >> ../readme.md
echo 'A cross-platform programmable network tool.' >> ../readme.md
echo '' >> ../readme.md
echo '# Sponsor' >> ../readme.md
echo '**❤️  [Shiliew - A network app designed for those who value their time](https://www.txthinking.com/shiliew.html)**' >> ../readme.md

cat getting-started.md >> ../readme.md
cat gui.md >> ../readme.md
cat resources.md >> ../readme.md

echo '# CLI Documentation' >> ../readme.md
jb '$1`brook mdpage`.split("\n").filter(v=>!v.startsWith("[")).join("\n").replace("```\n```", "```\nbrook --help\n```").split("\n").forEach(v=> echo(v.startsWith("**") && !v.startsWith("**Usage") ? "- "+v : v))' >> ../readme.md

cat example.md >> ../readme.md

markdown ../readme.md ./index.html

