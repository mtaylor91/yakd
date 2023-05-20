
.PHONY: all
all: node-controller

.PHONY: node-controller
node-controller:
	$(MAKE) -C node-controller

.PHONY: run
run: node-controller
	./node-controller/bin/node-controller run
