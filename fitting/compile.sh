#!/bin/bash
echo -e "Input cpp name:\c"
read hp
g++ -c $hp.cpp -o $hp.o \
 -I /usr/include \
 -L /usr/lib/x86_64-linux-gnu \
 -lopencv_core -lm
g++ -c $hp.c -o c$hp.o
gcc -o lib$hp.so c$hp.o $hp.o \
-lm -lopencv_core -lstdc++
