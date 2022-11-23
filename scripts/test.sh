#!/bin/sh
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
go test "$@" ./...
