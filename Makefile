.PHONY: install sync build clean test

LOCAL_BIN := $(shell go env GOPATH)/bin/gentle-ai
GIT_ROOT := $(shell git rev-parse --show-toplevel)

install:
	go install

build:
	go build -o $(LOCAL_BIN)

sync:
	$(LOCAL_BIN) sync

apply:
	$(LOCAL_BIN) apply

verify:
	$(LOCAL_BIN) verify

clean:
	rm -f $(LOCAL_BIN)
