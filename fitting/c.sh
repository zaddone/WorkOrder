#!/bin/bash
gcc -c cfitting.c -o cfitting.o \
 -I /usr/include \
 -L /usr/lib/x86_64-linux-gnu \
 -lopencv_core
ar r libcfitting.a cfitting.o
