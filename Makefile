.PHONY: run

build:
	@go build -o bin/vikingx
run: build 
	@bin/vikingx
test:
	@go test -v ./...
