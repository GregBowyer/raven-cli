GO ?= go
GOPATH := $(CURDIR)
export GOPATH

default: clean build

clean:
	$(GO) clean

build:
	$(GO) build 

fmt:
	$(GO) fmt

.PHONY: fmt build clean 
