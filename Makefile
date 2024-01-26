all: test build

build: 
	go build -o bin/ipinfo-server .

test:
	go test -v ./...

clean:
	rm -rf bin
