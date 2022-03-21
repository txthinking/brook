#!/bin/bash

# nami install markdown

markdown ../README.md index.html
gsed -i 's/README_ZH.md/index_zh.html/' index.html
markdown ../README_ZH.md index_zh.html
gsed -i 's/README.md/index.html/' index_zh.html
