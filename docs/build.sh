#!/bin/bash

# nami install markdown

markdown ../README.md index.html
gsed -i 's/README_ZH\.md/index_zh.html/' index.html
gsed -i 's/\.md/.html/' index.html

markdown ../README_ZH.md index_zh.html
gsed -i 's/README\.md/index.html/' index_zh.html
gsed -i 's/\.md/.html/' index_zh.html

rm -rf protocol
mkdir protocol
markdown ../protocol/brook-link-protocol.md protocol/brook-link-protocol.html
markdown ../protocol/brook-server-protocol.md protocol/brook-server-protocol.html
markdown ../protocol/brook-wsserver-protocol.md protocol/brook-wsserver-protocol.html
markdown ../protocol/brook-wssserver-protocol.md protocol/brook-wssserver-protocol.html
markdown ../protocol/withoutbrookprotocol-protocol.md protocol/withoutbrookprotocol-protocol.html
