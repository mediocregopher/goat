GOBIN=/bin/go

all: build

build:
	GOPATH=$(shell pwd):$(GOPATH) GOOS=linux GOARCH=386 $(GOBIN) build -o bin/goat_linux_386 src/goat/main/goat.go
	GOPATH=$(shell pwd):$(GOPATH) GOOS=linux GOARCH=amd64 $(GOBIN) build -o bin/goat_linux_amd64 src/goat/main/goat.go
