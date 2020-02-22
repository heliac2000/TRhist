##
## Makefile for trhist
##

SHELL = /bin/bash

SRCS = trhist.go approx.go lcp.go rmq.go util.go
TARGET = trhist
GCFLAGS = -gcflags='-B'
LDFLAGS = -ldflags='-s'
RELEASE_DIR = ../Release_V1

$(TARGET) b build: $(SRCS)
	@go build $(GCFLAGS) $(LDFLAGS) -o $(TARGET) .

test:
	@go test -v ./...

init_module:
#	@go mod init $(TARGET)
	@go mod init $(shell basename $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST)))))

tags:
	@find -name '*.go' | xargs etags

release rel:
	@test -d "$(RELEASE_DIR)" || mkdir "$(RELEASE_DIR)"
	@git archive --worktree-attributes --format=tar HEAD | tar -C "$(RELEASE_DIR)" -xf -

clean:
	@rm -f $(TARGET)
