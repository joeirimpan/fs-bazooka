BIN := $(shell basename $$PWD)

HASH := $(shell git rev-parse --short HEAD)
COMMIT_DATE := $(shell git show -s --format=%ci ${HASH})
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
VERSION := ${HASH} (${COMMIT_DATE})

build:
	@go build -o ${BIN} -ldflags="-X 'main.buildVersion=${VERSION}' -X 'main.buildDate=${BUILD_DATE}'"
	$(info Build successful. Current build version: $(VERSION))
.PHONY: build

test:
	@go test
.PHONY: test

run: build
	@./${BIN}
.PHONY: run

clean:
	@go clean
	-@rm -f ${BIN}
.PHONY: clean
