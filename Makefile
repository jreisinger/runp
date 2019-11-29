test:
	GO111MODULE=on go test ./...

build: test
	GO111MODULE=on go build

install: test
	GO111MODULE=on go install

