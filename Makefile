CC = gcc
CFLAGS = -Wall -Wextra -Werror
ALL_CFLAGS = $(CFLAGS) -Ilibyaml/install/include -Ilua/install/include
YL_LDFLAGS = -Llibyaml/install/lib -Llua/install/lib
YL_LDLIBS = -llua -lyaml -lm -largp

.PHONY: all
all: build/main.out

build/main.out: build/environment.o build/error.o build/event.o build/executor.o build/lua_helpers.o build/main.o build/parser.o build/render.o build/test.o
	$(CC) $(ALL_CFLAGS) $^ $(YL_LDFLAGS) $(YL_LDLIBS) -o $@

build/%.o: %.c
	mkdir -p build
	$(CC) $(ALL_CFLAGS) -c $< -o $@

deps.mk: *.c *.h deps.sh
	./deps.sh >$@

include deps.mk
include lua.mk
include libyaml.mk

.PHONY: clean
clean:
	rm -rf lua libyaml build
