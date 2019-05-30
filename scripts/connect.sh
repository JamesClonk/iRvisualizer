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
echo "connecting to iRvisualizer database ..."
psql ${IRVISUALIZER_DB_URI}
