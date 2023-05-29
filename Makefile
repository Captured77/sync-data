.PHONY: build install run

BINARY="sync-data"

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}

install:
	mv ./bin/${BINARY} /usr/local/bin

run:
	nohup ${BINARY}  >/dev/null 2>&1 & 

help:
	@echo "make build: build go executed file"
	@echo "make install: mv executed file to usr/local/bin"
	@echo "make run: run sync-data in the backgroud" 