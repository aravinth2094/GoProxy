.SILENT: clean deps test build

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=${GOCMD} test
GOGET=$(GOCMD) get
BINARY_NAME=GoProxy
BUILD_FILE=main.go
FLAGS=-ldflags "-s -w"
BUILD_DIRECTORY=dist

all: clean deps test build
test:
	${GOTEST} -v
build:
	echo "Creating build directory ${BUILD_DIRECTORY}..."
	mkdir ${BUILD_DIRECTORY}
	echo "Building darwin amd64..."
	GOOS=darwin GOARCH=amd64 ${GOBUILD} -o ${BUILD_DIRECTORY}/${BINARY_NAME}_darwin-amd64.out ${FLAGS} ${BUILD_FILE}
	echo "Building windows amd64..."
	GOOS=windows GOARCH=amd64 ${GOBUILD} -o ${BUILD_DIRECTORY}/${BINARY_NAME}_win-amd64.exe ${FLAGS} ${BUILD_FILE}
	echo "Building windows x86..."
	GOOS=windows GOARCH=386 ${GOBUILD} -o ${BUILD_DIRECTORY}/${BINARY_NAME}_win-x86.exe ${FLAGS} ${BUILD_FILE}
	echo "Building linux amd64..."
	GOOS=linux GOARCH=amd64 ${GOBUILD} -o ${BUILD_DIRECTORY}/${BINARY_NAME}_linux-amd64.out ${FLAGS} ${BUILD_FILE}
clean:
	echo "Cleaning ${BUILD_DIRECTORY}..."
	rm -rf ${BUILD_DIRECTORY}
deps:
	echo "Fetching dependencies..."
	go get -u