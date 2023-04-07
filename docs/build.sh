#!/bin/bash

echo '# Brook CLI Documentation' > cli.md
brook mdpage >> cli.md
echo '<!--SIDEBAR-->' >> cli.md
echo '<!--G-R3M673HK5V-->' >> cli.md
markdown cli.md
markdown gui.md
markdown gui-zh.md
markdown example.md
markdown example-zh.md
markdown diagram.md diagram.html

echo '# Documentation' > ../readme.md
echo >> ../readme.md
echo '👉 [**Documentation**](https://txthinking.github.io/brook/)' >> ../readme.md
cat getting-started.md >> ../readme.md
echo '# Brook CLI Documentation' >> ../readme.md
brook mdpage >> ../readme.md
cat gui.md >> ../readme.md
cat diagram.md >> ../readme.md
mdtoc ../readme.md > toc.md

echo '# Brook' > ../readme.md
echo >> ../readme.md
echo 'A cross-platform network tool designed for developers.' >> ../readme.md
echo >> ../readme.md
echo '[❤️  *A txthinking project*](https://www.txthinking.com)' >> ../readme.md
echo >> ../readme.md
cat toc.md >> ../readme.md
rm -rf toc.md
echo '# Documentation' >> ../readme.md
echo >> ../readme.md
echo '👉 [**Documentation**](https://txthinking.github.io/brook/)' >> ../readme.md
cat getting-started.md >> ../readme.md
echo '# Brook CLI Documentation' >> ../readme.md
brook mdpage >> ../readme.md
cat gui.md >> ../readme.md
cat diagram.md >> ../readme.md
