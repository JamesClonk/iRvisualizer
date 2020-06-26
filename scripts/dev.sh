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
source ~/.config/ircollector_db.conf || true

# =============================================================================================
echo "developing irvisualizer ..."
killall gin-bin || true
killall irvisualizer || true
rm -f gin-bin || true
#gin --all run main.go

rm -f irvisualizer || true
GOARCH=amd64 GOOS=linux go build -i -o irvisualizer
./irvisualizer
