SOURCES := $(shell find . -name '*.go')
INSTALL_DIR ?= /usr/local/bin

build: gentee

install: build
	cp gentee ${INSTALL_DIR}/gentee

clean:
	rm gentee

gentee: $(SOURCES)
	go build -o gentee ./cli

