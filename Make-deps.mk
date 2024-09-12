build/environment.o: environment.h lua/install
build/error.o: error.h libyaml/install lua/install
build/event.o: error.h event.h executor.h libyaml/install lua/install parser.h render.h
build/executor.o: error.h event.h executor.h libyaml/install lua/install lua_helpers.h parser.h render.h
build/lua_helpers.o: error.h event.h libyaml/install lua/install lua_helpers.h
build/main.o: environment.h error.h event.h executor.h libyaml/install lua/install parser.h render.h test.h
build/parser.o: error.h libyaml/install lua/install parser.h
build/render.o: error.h event.h executor.h libyaml/install lua/install lua_helpers.h parser.h render.h
build/test.o: error.h event.h executor.h libyaml/install lua/install parser.h render.h test.h
build/ylt.o: libyaml/install lua/install ylt.h
build/main.out: build/environment.o build/error.o build/event.o build/executor.o build/lua_helpers.o build/main.o build/parser.o build/render.o build/test.o build/ylt.o
