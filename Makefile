MODULES := $(wildcard modules/*/.)
.PHONY: all get build $(MODULES)

all: build $(MODULES)

deps: get $(MODULES)

get:
	go get

build:
	go build .

$(MODULES):
	$(MAKE) -C $@ $(MAKECMDGOALS)

