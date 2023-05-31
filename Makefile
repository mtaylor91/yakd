
GO_SRCS := $(shell find . -name '*.go')

.PHONY: all
all: bin/yakd

.PHONY: clean
clean:
	rm -rf bin
	rm -rf build

.PHONY: test
test:
	go test ./...

bin/yakd: $(GO_SRCS)
	@mkdir -p $(@D)
	go build -o $@

.PHONY: run
run: node-controller
	./node-controller/bin/node-controller run
