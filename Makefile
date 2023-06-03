
GO_BUILD_SRCS := $(shell find cmd/build -name '*.go')
GO_PKG_SRCS := $(shell find pkg -name '*.go')

.PHONY: all
all: bin/yakd-build

.PHONY: clean
clean: clean-bin clean-build

.PHONY: clean-bin
clean-bin:
	rm -rf bin

.PHONY: clean-build
clean-build:
	rm -rf build

.PHONY: test
test:
	go test ./...

bin/yakd-build: $(GO_BUILD_SRCS) $(GO_PKG_SRCS)
	@mkdir -p $(@D)
	go build -o $@ ./cmd/build

.PHONY: run
run: node-controller
	./node-controller/bin/node-controller run
