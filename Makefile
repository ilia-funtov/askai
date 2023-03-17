build:
	go build -o bin/askai

test:
	go test

lint:
	golangci-lint run --enable-all

run: build
	bin/askai

clean:
	go clean
	rm -rf bin