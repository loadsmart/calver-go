GOFILES = $(shell find . -name '*.go' -not -path './vendor/*' -not -path '\.*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

test:
	@go test -v $(GOPACKAGES)

lint:
	@golangci-lint run -v --fix
