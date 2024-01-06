#!/usr/bin/env bash
set -euo pipefail

./main.out -i testcases/call.yaml -t
# ./main.out -i testcases/eval.yaml -t
# ./main.out -i testcases/for.yaml -t
./main.out -i testcases/formatting.yaml -t
# ./main.out -i testcases/if.yaml -t
