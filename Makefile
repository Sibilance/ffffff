lua:
	mkdir -p lua
	curl https://www.lua.org/ftp/lua-5.4.6.tar.gz | tar xvzC lua --strip-components=1

lua/install: lua
	cd lua && make all local

lua/install/include/lua.h: lua/install
lua/install/include/lauxlib.h: lua/install
lua/install/include/lualib.h: lua/install
lua/install/lib/liblua.a: lua/install

libyaml:
	mkdir -p libyaml
	curl https://pyyaml.org/download/libyaml/yaml-0.2.5.tar.gz | tar xvzC libyaml --strip-components=1

libyaml/install: libyaml
	cd libyaml && ./configure --prefix="${CURDIR}/libyaml/install"
	cd libyaml && make all install

libyaml/install/include/yaml.h: libyaml/install
libyaml/install/lib/libyaml.a: libyaml/install

.PHONY: clean
clean:
	rm -rf lua libyaml
