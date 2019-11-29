test:
	GO111MODULE=on go test ./...

build: test
	GO111MODULE=on go build

install: test
	GO111MODULE=on go install

PLATFORMS := linux/amd64 darwin/amd64 linux/arm

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

release: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o 'runp-$(os)-$(arch)' main.go

release: $(PLATFORMS)
