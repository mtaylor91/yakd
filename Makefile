
GO_SRCS := $(shell find . -name '*.go')

.PHONY: all
all: bin/yakd

.PHONY: clean
clean:
	rm -rf bin
	rm -f stage1.tar.gz
	rm -f yakd.qcow2
	rm -f yakd.qcow2.raw

bin/yakd: $(GO_SRCS)
	@mkdir -p $(@D)
	go build -o $@

.PHONY: run
run: node-controller
	./node-controller/bin/node-controller run
