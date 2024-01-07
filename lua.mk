lua:
	mkdir -p $@
	curl https://www.lua.org/ftp/lua-5.4.6.tar.gz | tar xvzC $@ --strip-components=1

lua/install: lua
	cd lua && make all local
