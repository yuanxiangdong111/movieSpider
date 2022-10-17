GOCMD			:=$(shell which go)
GOBUILD			:=$(GOCMD) build


IMPORT_PATH		:=btspyder/cmd
BUILD_TIME		:=$(shell date "+%F %T")
COMMIT_ID       :=$(shell git rev-parse HEAD)
GO_VERSION      :=$(shell $(GOCMD) version)
#VERSION			:=$(shell git describe --tags)
VERSION			:=v1.0
BUILD_USER		:=$(shell whoami)
FLAG			:="-X '${IMPORT_PATH}.buildTime=${BUILD_TIME}' -X '${IMPORT_PATH}.commitID=${COMMIT_ID}' -X '${IMPORT_PATH}.goVersion=${GO_VERSION}' -X '${IMPORT_PATH}.goVersion=${GO_VERSION}' -X '${IMPORT_PATH}.Version=${VERSION}' -X '${IMPORT_PATH}.buildUser=${BUILD_USER}'"

BINARY_DIR=bin/movieSpider
BINARY_NAME:=movieSpider



build:
	CGO_ENABLED=1 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)
# linux
build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-linux

#mac
build-darwin:
	CGO_ENABLED=1 GOOS=darwin $(GOBUILD) -ldflags $(FLAG) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin

# windows
build-win:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG)  -o $(BINARY_DIR)/$(BINARY_NAME)-win.exe

# 全平台
build-all:
	make build-linux
	make build-darwin
	make build-win
	cd bin&&tar zcf ${BINARY_NAME}.tgz ${BINARY_NAME}

#docker
build-image:
	make build
	docker build -t harbor.dlab.cn/public/btspyder:$(VERSION) .

#docker
push-image:
	docker push harbor.dlab.cn/public/btspyder:$(VERSION)
