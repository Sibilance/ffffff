#!/bin/sh
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
mkdir -p lua
curl https://www.lua.org/ftp/lua-5.4.6.tar.gz | tar xvzC lua --strip-components=1
