#!/bin/bash
. build.sh
cd ../
git add -A .
git commit -a -m "publish site"
git push
