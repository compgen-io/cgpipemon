#GOPATH=$(shell pwd)/vendor:$(shell pwd)
#GOBIN=$(shell pwd)/bin
#GOFILES=$(wildcard *.go)
#GONAME=$(shell basename "$(PWD)")

build:
	@echo "Bundling resources with go-bindata"
	@go-bindata -o config/bindata.go -pkg config db/schema.sql
	@echo "Building $(GOFILES) to ./bin"
#	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/cgpipemon-server cmd/server.go
	go build -o bin/cgpipemon-server cmd/server.go

clean:
	@echo "Cleaning"
	go clean
	rm config/bindata.go

.PHONY: build clean
