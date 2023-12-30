CC = gcc
CFLAGS = -Wall -Wextra -Werror
ALL_CFLAGS = $(CFLAGS) -Ilibyaml/install/include -Ilua/install/include
YL_LDFLAGS = -Llibyaml/install/lib -Llua/install/lib
YL_LDLIBS = -llua -lyaml -largp

.PHONY: all
all: main.out

main.out: environment.o error.o event.o executor.o main.o parser.o render.o test.o
	$(CC) $(ALL_CFLAGS) $^ $(YL_LDFLAGS) $(YL_LDLIBS) -o main.out

environment.o: lua/install
error.o: libyaml/install
event.o: lua/install libyaml/install
executor.o: lua/install libyaml/install
main.o: lua/install libyaml/install
parser.o: libyaml/install
render.o: lua/install libyaml/install

environment.c: environment.h
error.c: error.h render.h
event.c: event.h
executor.c: executor.h
main.c: environment.h executor.h parser.h render.h test.h
parser.c: parser.h
render.c: render.h
test.c: test.h

parser.h: error.h
event.h: error.h
executor.h: event.h parser.h
render.h: executor.h
test.h: event.h executor.h parser.h

%.o: %.c 
	$(CC) $(ALL_CFLAGS) -c $<

%.h %.c:  # Propagate file save timestamps onto dependees.
	touch $@

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
