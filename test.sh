#!/usr/bin/env bash
set -euo pipefail

build/main.out -i testcases/call.yaml -t
# build/main.out -i testcases/eval.yaml -t
# build/main.out -i testcases/for.yaml -t
build/main.out -i testcases/formatting.yaml -t
# build/main.out -i testcases/if.yaml -t
