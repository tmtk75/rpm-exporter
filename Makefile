VERSION := $(shell git describe --tags --abbrev=0)
VERSION_LONG := $(shell git describe --tags)
VAR_VERSION := main.Version

LDFLAGS := -ldflags "-X $(VAR_VERSION)=$(VERSION) \
	-X $(VAR_VERSION)Long=$(VERSION_LONG)"

rpm-exporter: main.go
	make build NAME=rpm-exporter

rpm-exporter_linux_amd64:
	GOARCH=amd64 GOOS=linux make build NAME=rpm-exporter_linux_amd64

NAME := rpm-exporter
.PHONY: build
build:
	go build $(LDFLAGS) -o $(NAME) main.go

