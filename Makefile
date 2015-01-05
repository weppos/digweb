name = digweb

test:
	go test -v ./...

build:
	go build -o $(name)

get-deps:
	go get -d -v ./...
