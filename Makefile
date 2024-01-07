CC = gcc
CFLAGS = -Wall -Wextra -Werror
ALL_CFLAGS = $(CFLAGS) -Ilibyaml/install/include -Ilua/install/include
YL_LDFLAGS = -Llibyaml/install/lib -Llua/install/lib
YL_LDLIBS = -llua -lyaml -lm -largp

.PHONY: all
all: build/main.out

build/main.out:
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
