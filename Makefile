BUILD_DIR ?= build

.PHONY: build package clean

build:
	mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/main *.go

package:
	mkdir -p ${BUILD_DIR}
	cp support/wrapper/index.js ${BUILD_DIR}/index.js
	cp support/wrapper/byline.js ${BUILD_DIR}/byline.js
	zip -j ${BUILD_DIR}/archive.zip ${BUILD_DIR}/index.js ${BUILD_DIR}/byline.js ${BUILD_DIR}/main

clean:
	rm -rf ${BUILD_DIR}

all: clean build package
