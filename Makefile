
GO_SRCS := $(shell find . -name '*.go')

.PHONY: all
all: bin/yakd

bin/yakd: $(GO_SRCS)
	@mkdir -p $(@D)
	go build -o $@

.PHONY: run
run: node-controller
	./node-controller/bin/node-controller run
