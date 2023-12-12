CC = gcc
CFLAGS = -Wall -Werror
ALL_CFLAGS = $(CFLAGS) -Ilibyaml/install/include -Ilua/install/include
YL_LDFLAGS = -Llibyaml/install/lib -Llua/install/lib
YL_LDLIBS = -llua -lyaml -largp

.PHONY: all
all: main.out

main.out: error.o main.o parser.o
	$(CC) $(ALL_CFLAGS) $^ $(YL_LDFLAGS) $(YL_LDLIBS) -o main.out

error.o: error.c libyaml/install
	$(CC) $(ALL_CFLAGS) -c error.c

main.o: main.c lua/install libyaml/install
	$(CC) $(ALL_CFLAGS) -c main.c

parser.o: parser.c libyaml/install
	$(CC) $(ALL_CFLAGS) -c parser.c

error.c: error.h
main.c: parser.h
parser.c: parser.h

parser.h: error.h

lua:
	mkdir -p lua
	curl https://www.lua.org/ftp/lua-5.4.6.tar.gz | tar xvzC lua --strip-components=1

lua/install: lua
	cd lua && make all local

libyaml:
	mkdir -p libyaml
	curl https://pyyaml.org/download/libyaml/yaml-0.2.5.tar.gz | tar xvzC libyaml --strip-components=1

libyaml/install: libyaml
	cd libyaml && ./configure --prefix="${CURDIR}/libyaml/install"
	cd libyaml && make all install

.PHONY: clean
clean:
	rm -rf lua libyaml *.o *.out
