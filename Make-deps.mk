build/error.o: error.h libyaml/install
build/event.o: error.h event.h libyaml/install
build/executor.o: error.h event.h executor.h libyaml/install parser.h
build/main.o: error.h event.h executor.h libyaml/install parser.h
build/parser.o: error.h libyaml/install parser.h
build/main.out: build/error.o build/event.o build/executor.o build/main.o build/parser.o
