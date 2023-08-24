.PHONY: all

build:
	@go build -o bin/vikingx
server:
	cd frontend && npm run build

run: build server
	@bin/vikingx
test:
	@go test -v ./...
