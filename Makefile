UNAME_S := $(shell uname -s)

# fire up golangci-lint to concurrently run several static analysis tools at once
.PHONY: lint
lint:
	GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	# recursively run golangci-lint on all files in this directory, skipping packages in vendor
	golangci-lint run --enable-all --skip-dirs vendor

# build everything in this directory into a single binary in bin-directory
.PHONY: build
build:
ifeq ($(OS),Windows_NT)
	GO111MODULE=on GOOS=windows GOARCH=386 go build -o bin/main.exe
else
	GO111MODULE=on go build -o bin/main
endif

# build docker image
.PHONY: docker
docker: build
ifeq ($(UNAME_S),Linux)
	sudo docker build -t trivago/monitoring:edgecast-v1 .
else
	docker build -t trivago/monitoring:edgecast-v1 .
endif
