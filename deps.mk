environment.o: environment.h lua/install
error.o: error.h libyaml/install
event.o: error.h event.h executor.h libyaml/install lua/install parser.h render.h
executor.o: error.h event.h executor.h libyaml/install lua/install parser.h
main.o: environment.h error.h event.h executor.h libyaml/install lua/install parser.h render.h test.h
parser.o: error.h libyaml/install parser.h
render.o: error.h event.h executor.h libyaml/install lua/install parser.h render.h
test.o: error.h event.h executor.h libyaml/install lua/install parser.h test.h
