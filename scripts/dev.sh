#!/bin/bash

# fail on error
set -e

# =============================================================================================
if [[ "$(basename $PWD)" == "scripts" ]]; then
    cd ..
fi
echo $PWD

# =============================================================================================
source .env
source ~/.config/ir.conf || true

# =============================================================================================
echo "developing iRvisualizer ..."
killall gin-bin || true
killall iRvisualizer || true
rm -f gin-bin || true
#gin --all run main.go

rm -f iRvisualizer || true
GOARCH=amd64 GOOS=linux go build -i -o iRvisualizer
./iRvisualizer
