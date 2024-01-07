libyaml:
	mkdir -p $@
	curl https://pyyaml.org/download/libyaml/yaml-0.2.5.tar.gz | tar xvzC $@ --strip-components=1

libyaml/Makefile: libyaml
	cd libyaml && ./configure --prefix="${CURDIR}/libyaml/install"
	touch libyaml/Makefile

libyaml/install: libyaml/Makefile
	cd libyaml && make all install
	touch libyaml/Makefile libyaml/install
