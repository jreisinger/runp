test:
	GO111MODULE=on go test ./...

install: test
	GO111MODULE=on go install

PLATFORMS := linux/amd64 linux/arm darwin/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

release: test $(PLATFORMS)

$(PLATFORMS):
	GO111MODULE=on GOOS=$(os) GOARCH=$(arch) go build -o 'runp-$(os)-$(arch)' main.go
