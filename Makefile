.PHONY: build

build: generate
	go build .

.PHONY: generate
generate:
	go mod vendor
	go generate ./...

.PHONY: test
test:
	go test .