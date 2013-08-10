GOBIN=/bin/go

all: build

build:
	GOPATH=$(shell pwd):$(GOPATH) $(GOBIN) build -o bin/goat src/main/goat.go
