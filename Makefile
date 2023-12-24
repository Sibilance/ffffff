CC = gcc
CFLAGS = -Wall -Wextra -Werror
ALL_CFLAGS = $(CFLAGS) -Ilibyaml/install/include -Ilua/install/include
YL_LDFLAGS = -Llibyaml/install/lib -Llua/install/lib
YL_LDLIBS = -llua -lyaml -largp

.PHONY: all
all: main.out

main.out: error.o executor.o main.o parser.o producer.o test.o
	$(CC) $(ALL_CFLAGS) $^ $(YL_LDFLAGS) $(YL_LDLIBS) -o main.out

error.o: libyaml/install
executor.o: lua/install libyaml/install
main.o: lua/install libyaml/install
parser.o: libyaml/install

error.c: error.h
executor.c: executor.h producer.h
main.c: executor.h parser.h test.h
parser.c: parser.h
producer.c: producer.h
test.c: test.h

parser.h: error.h
executor.h: parser.h
producer.h: executor.h
test.h: executor.h parser.h

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

libyaml/install: libyaml
	cd libyaml && ./configure --prefix="${CURDIR}/libyaml/install"
	cd libyaml && make all install

.PHONY: clean
clean:
	rm -rf lua libyaml *.o *.out
