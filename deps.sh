#!/usr/bin/env bash

declare -A depcache

function dependencies() {
    local deps=${depcache[$1]}
    if [ -z "$deps" ]; then
        deps=$(sed -nre 's/^#include +"([^"]+)"$/\1/p' "$1")
        depcache[$1]=$deps
    fi
    for file in $deps; do
        case "$file" in 
            "lua.h"|"lauxlib.h"|"lualib.h")
                echo "lua/install"
                ;;
            "yaml.h")
                echo "libyaml/install"
                ;;
            *)
                echo "$file"
                dependencies "$file"
                ;;
        esac
    done
}

for file in *.c; do
    echo "${file/%.c/.o}:" $(dependencies "$file" | sort | uniq)
done
