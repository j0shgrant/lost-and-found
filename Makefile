LAF_BUILD_VERSION := $(shell git describe --tags)
REPO_PATH := github.com/j0shgrant/lost-and-found

default: clean install

build: build-laf

install: install-laf

build-laf:
	@go build -ldflags="-X '$(REPO_PATH)/laf/cmd.version=$(LAF_BUILD_VERSION)'" -o "./bin/" "./laf"

install-laf:
	@go install -ldflags="-X '$(REPO_PATH)/laf/cmd.version=$(LAF_BUILD_VERSION)'" "./laf"

clean:
	@if [ -d "./bin" ]; then rm -rf "./bin"; fi
