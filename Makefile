SOURCES := $(shell find . -name '*.go')

gentee: $(SOURCES)
	go build -o gentee ./cli
