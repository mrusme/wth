MODULES := $(wildcard modules/*/.)
.PHONY: all build $(MODULES)

all: build $(MODULES)

build:
	go build .

$(MODULES):
	$(MAKE) -C $@

