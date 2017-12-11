.PHONY: clean current vendor test
BUILD_FLAGS=GO15VENDOREXPERIMENT=1 GORACE="halt_on_error=0" GOGC=off

current:
	@$(BUILD_FLAGS) go build

clean:
	@rm -f ./edgecast
	@rm -f ./dist/edgecast_*.zip
	@go clean

vendor:
	@go get -u github.com/FiloSottile/gvt
	@gvt update -all -precaire

test:
	@$(BUILD_FLAGS) go test -cover -v -timeout 10s -race $$(go list ./...|grep -v vendor)

.DEFAULT_GOAL := current