.PHONY: build test run lint lint-all clean vulncheck clean

build:
	go build -o bin/askai

install:
	go install

test:
	go test

run: build
	bin/askai

lint:
	golangci-lint run --enable revive,gomnd,goimports,gosec,funlen,cyclop,gocognit,wrapcheck,errorlint

lint-all:
	golangci-lint run --enable-all

vulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck .

clean:
	go clean
	rm -rf bin