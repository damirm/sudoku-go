main: main.go
	./build.sh

.PHONY: default
default: main

run: main
	./main

test:
	go test -race -v
