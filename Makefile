CC = gcc
CFLAGS = -Wall -Wextra -Werror
ALL_CFLAGS = $(CFLAGS) -Ilibyaml/install/include -Ilua/install/include
YL_LDFLAGS = -Llibyaml/install/lib -Llua/install/lib
YL_LDLIBS = -llua -lyaml -largp

.PHONY: all
all: main.out

main.out: environment.o error.o event.o executor.o main.o parser.o render.o test.o
	$(CC) $(ALL_CFLAGS) $^ $(YL_LDFLAGS) $(YL_LDLIBS) -o main.out

%.o: %.c 
	$(CC) $(ALL_CFLAGS) -c $<

deps.mk: *.c *.h
	./deps.sh >deps.mk
include deps.mk

lua:
	mkdir -p lua
	curl https://www.lua.org/ftp/lua-5.4.6.tar.gz | tar xvzC lua --strip-components=1

lua/install: lua
	cd lua && make all local

libyaml:
	mkdir -p libyaml
	curl https://pyyaml.org/download/libyaml/yaml-0.2.5.tar.gz | tar xvzC libyaml --strip-components=1

libyaml/Makefile: libyaml
	cd libyaml && ./configure --prefix="${CURDIR}/libyaml/install"
	touch libyaml/Makefile

libyaml/install: libyaml/Makefile
	cd libyaml && make all install
	touch libyaml/Makefile libyaml/install

.PHONY: clean
clean:
	rm -rf lua libyaml *.o *.out
