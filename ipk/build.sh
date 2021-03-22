#!/bin/bash

cd control
tar -czf control.tar.gz *
mv control.tar.gz ../
cd ..

cd data
tar -czf data.tar.gz *
mv data.tar.gz ../
cd ..

tar -czf brook.ipk control.tar.gz data.tar.gz debian-binary

rm *.tar.gz
