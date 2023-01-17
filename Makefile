SOURCES := $(shell find . -name '*.go')
INSTALL_DIR ?= /usr/local/bin

NVIM_SYNTAX_DIR ?=$(HOME)/.config/nvim/syntax

build: gentee

install: build
	cp gentee ${INSTALL_DIR}/gentee

install-nvim:
	mkdir -p ${NVIM_SYNTAX_DIR}
	cp contrib/gentee.vim ${NVIM_SYNTAX_DIR}/gentee.vim

clean:
	rm gentee

gentee: $(SOURCES)
	go build -o gentee ./cli

