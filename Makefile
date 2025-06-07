BINARY=bin/dynago
INSTALL_BIN=/usr/local/bin/dynago
CONFIG_SRC=configs/dynago.yml
CONFIG_DST=/etc/dynago/dynago.yml
SYSTEMD_UNIT=/etc/systemd/system/dynago.service

VERSION := 0.1.1
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD)

LDFLAGS := -X 'main.Version=$(VERSION)' \
           -X 'main.BuildTime=$(BUILD_TIME)' \
           -X 'main.GitCommit=$(GIT_COMMIT)'

.PHONY: all build install clean fmt test

all: build

build:
	mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/dynago

fmt:
	go fmt ./...

test:
	go test ./...

install: build
	install -Dm755 $(BINARY) $(INSTALL_BIN)
	install -d /etc/dynago
	install -Dm644 $(CONFIG_SRC) $(CONFIG_DST)
	install -Dm644 dynago.service $(SYSTEMD_UNIT)
	systemctl daemon-reload
	systemctl enable dynago.service

clean:
	rm -rf bin