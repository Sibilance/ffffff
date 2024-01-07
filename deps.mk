build/environment.o: environment.h lua/install
build/error.o: error.h libyaml/install
build/event.o: error.h event.h executor.h libyaml/install lua/install parser.h render.h
build/executor.o: error.h event.h executor.h libyaml/install lua/install parser.h render.h
build/main.o: environment.h error.h event.h executor.h libyaml/install lua/install parser.h render.h test.h
build/parser.o: error.h libyaml/install parser.h
build/render.o: error.h event.h executor.h libyaml/install lua/install parser.h render.h
build/test.o: error.h event.h executor.h libyaml/install lua/install parser.h render.h test.h
