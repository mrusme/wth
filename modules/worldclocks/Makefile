MODNAME=$(shell basename $(CURDIR))

build:
	go build -buildmode=plugin -o $(MODNAME).so .

deps:
	go get
