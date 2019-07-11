VERSION := $(shell git describe --tags --abbrev=0)
VERSION_LONG := $(shell git describe --tags)
VAR_VERSION := main.Version

LDFLAGS := -ldflags "-X $(VAR_VERSION)=$(VERSION) \
	-X $(VAR_VERSION)Long=$(VERSION_LONG)"

build: ./rpm-exporter

rpm-exporter: main.go
	go build $(LDFLAGS) -o rpm-exporter main.go
