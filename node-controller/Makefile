
GO_SRCS := $(shell find . -name '*.go')

.PHONY: all

all: bin/node-controller

bin/node-controller: $(GO_SRCS)
	@mkdir -p $(@D)
	go build -o $@
