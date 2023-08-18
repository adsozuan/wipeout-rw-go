
## build: build application cmd/wipeout into bin/
.PHONY: build
build:
	go build -ldflags='-s' -o=./bin/wipeout.exe ./cmd/wipeout
