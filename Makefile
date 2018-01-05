SHELL := /bin/bash
APP="pgme"
PROJECT?=github.com/chhibber/pgme
PORT?=9101

RELEASE?=$(shell bash ./build/ver.sh)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
CONTAINER_IMAGE?=chhibber/${APP}

GOOS?=linux
GOARCH?=amd64

.PHONY: clean
clean:
	rm -f ${APP}


.PHONY: docker-build
docker-build:
	docker build -t ${IMAGE_REPO_NAME}:${IMAGE_TAG} .


.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w \
		-X main.Release=${RELEASE} \
		-X main.Commit=${COMMIT} \
		-X main.BuildTime=${BUILD_TIME}" \
		-o ${APP}

build-mac: clean
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
			-ldflags "-s -w \
			-X main.Release=${RELEASE} \
			-X main.Commit=${COMMIT} \
			-X main.BuildTime=${BUILD_TIME}" \
			-o ${APP}
