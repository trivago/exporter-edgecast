# dynamically fetch path to executables
GO_BIN := $(GOPATH)/bin
GOMETALINTER := $(GO_BIN)/gometalinter

# in case gometalinter is not installed already => clone it and install it
$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

# fire up gometalinter to concurrently run several static analysis tools at once
.PHONY: lint
lint: $(GOMETALINTER)
	# recursively run gometalinter on all files in this directory, skipping packages in vendor
	gometalinter ./... --vendor --disable=gotype

.PHONY: dependencies
dependencies:
	dep ensure

# build everything in this directory into a single binary in bin-directory
.PHONY: build
build: dependencies
ifeq ($(OS),Windows_NT)
	go build -o bin/main.exe
else
	go build -o bin/main
endif

# build docker image
.PHONY: docker
docker: build
ifeq ($(OS),Windows_NT)
	docker build -t trivago/monitoring:edgecast-v1 .
else
	sudo docker build -t trivago/monitoring:edgecast-v1 .
endif