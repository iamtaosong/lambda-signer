BUILD_DIR ?= build

.PHONY: build config package clean

build:
	mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/main *.go

config:
	cp support/config.tpl ${BUILD_DIR}/config.json
	sed -i '' -e "s;__CONFIG_URL__;${CONFIG_URL};g" ${BUILD_DIR}/config.json

package:
	mkdir -p ${BUILD_DIR}
	cp support/wrapper/index.js ${BUILD_DIR}/index.js
	cp support/wrapper/byline.js ${BUILD_DIR}/byline.js
	zip -j ${BUILD_DIR}/archive.zip ${BUILD_DIR}/index.js ${BUILD_DIR}/byline.js ${BUILD_DIR}/config.json ${BUILD_DIR}/main

clean:
	rm -rf ${BUILD_DIR}

all: clean build config package
