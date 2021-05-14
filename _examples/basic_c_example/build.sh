#!/bin/bash
# archive here for static linking but can be buildmode=c-shared for dynamically loading
go build -buildmode=c-archive -o libunifiedpush.a ../../api_c/api.go 
gcc -Wall  -g main.c -Wno-unused-function -L . -lunifiedpush -pthread
