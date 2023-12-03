lua:
	scripts/install-lua.sh

lua/install: lua
	cd lua && make all local

.PHONY: clean
clean:
	rm -rf lua
