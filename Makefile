GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")

build:
	@echo "Bundling resources with go-bindata"
	@go-bindata -o config/bindata.go -pkg config db/schema.sql
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get .

install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

run:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES)

#watch:
#	@$(MAKE) restart &
#	@fswatch -o . -e 'bin/.*' | xargs -n1 -I{}  make restart

#restart: clear stop clean build start

# start:
#   @echo "Starting bin/$(GONAME)"
#   @./bin/$(GONAME) & echo $$! > $(PID)

# stop:
#   @echo "Stopping bin/$(GONAME) if it's running"
#   @-kill `[[ -f $(PID) ]] && cat $(PID)` 2>/dev/null || true

# clear:
#   @clear

clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
	rm bin/$(GONAME)
	rm config/bindata.go

.PHONY: build get install run  clean
