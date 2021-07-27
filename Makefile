.SILENT: clean deps test build

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=${GOCMD} test
GOGET=$(GOCMD) get
BINARY_NAME=GoProxy
BUILD_FILE=main.go
FLAGS=-ldflags "-s -w"
BUILD_DIRECTORY=dist
GORELEASERCMD=goreleaser
FLAGS=--skip-publish --snapshot --rm-dist

all: test build
test:
	${GOTEST} -v -cover
build:
	${GORELEASERCMD} ${FLAGS}