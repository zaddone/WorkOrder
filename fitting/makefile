TARG=fitting
CGOFILES=fitting.go
fitting.o:fitting.cpp
	g++ -I /usr/include -L /usr/lib/x86_64-linux-gnu -lm -lopencv_core -o $@ -c $(CGO_CFLAGS) $< $
cfitting.o:fitting.c
	g++  -fPIC -O2 -o $@ -c $(CGO_CFLAGS) $<
CGO_LDFLAGS+=-lstdc++
CGO_LDFLAGS+=-lopencv_core
$(elem)_fitting.so: fitting.cgo4.o fitting.o cfitting.o
	gcc $(_CGO_CFLAGS_$(GOARCH)) $(_CGO_LDFLAGS_$(GOOS)) -o $@ $^ $(CGO_LDFLAGS)
