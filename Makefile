CC = gcc
CFLAGS = -Wall -Wextra -Werror
ALL_CFLAGS = $(CFLAGS) -Ilibyaml/install/include -Ilua/install/include
YL_LDFLAGS = -Llibyaml/install/lib -Llua/install/lib
YL_LDLIBS = -llua -lyaml -lm -largp

.PHONY: all
all: build/main.out

build/main.out: build/environment.o build/error.o build/event.o build/executor.o build/main.o build/parser.o build/render.o build/test.o
	$(CC) $(ALL_CFLAGS) $^ $(YL_LDFLAGS) $(YL_LDLIBS) -o $@

build/%.o: %.c
	mkdir -p build
	$(CC) $(ALL_CFLAGS) -c $< -o $@

deps.mk: *.c *.h deps.sh
	./deps.sh >$@
include deps.mk

lua:
	mkdir -p $@
	curl https://www.lua.org/ftp/lua-5.4.6.tar.gz | tar xvzC $@ --strip-components=1

lua/install: lua
	cd lua && make all local

libyaml:
	mkdir -p $@
	curl https://pyyaml.org/download/libyaml/yaml-0.2.5.tar.gz | tar xvzC $@ --strip-components=1

libyaml/Makefile: libyaml
	cd libyaml && ./configure --prefix="${CURDIR}/libyaml/install"
	touch libyaml/Makefile

libyaml/install: libyaml/Makefile
	cd libyaml && make all install
	touch libyaml/Makefile libyaml/install

.PHONY: clean
clean:
	rm -rf lua libyaml build
