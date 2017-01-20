#!/bin/sh
set -e

./build.sh
git pull -r
git push
