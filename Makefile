.PHONY: all build modules

all: build modules

build:
	go build .

modules:
	@make -C modules/example
	@make -C modules/worldclocks

