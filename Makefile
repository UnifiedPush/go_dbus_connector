lint:
	golangci-lint run

c-static:
	go build -buildmode=c-archive -o bin/libunifiedpush.a ./api_c/
c-so:
	go build -buildmode=c-shared -o bin/libunifiedpush.so ./api_c/
test: c-so
	go test ./... #todo

all: c-static c-so
