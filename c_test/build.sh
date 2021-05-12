go build -buildmode=c-archive -o libunifiedpush.a ../api_c/api.go 
gcc -Wall -g main.c -L . -lunifiedpush -pthread

