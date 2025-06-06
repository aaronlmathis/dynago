BINARY=bin/dynago
INSTALL_BIN=/usr/local/bin/dynago
CONFIG_SRC=configs/dynago.yml
CONFIG_DST=/etc/dynago/dynago.yml
SYSTEMD_UNIT=/etc/systemd/system/dynago.service

.PHONY: all build install clean fmt test

all: build

build:
	mkdir -p bin
	go build -o $(BINARY) ./cmd/dynago

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