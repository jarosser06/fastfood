export GOPATH := ${PWD}/.gopath
export FFVERSION := 0.2.2
export INSTALLPRE := /usr/local

all: deps build

build:
		@echo "Building binary..."
		scripts/make.sh build

cross_compile:
		@echo "Building binaries..."
		go get github.com/mitchellh/gox
		@${GOPATH}/bin/gox -arch="amd64" -os="linux darwin windows" -output="bin/{{.OS}}_{{.Arch}}/fastfood" ./cmd/main/

package: clean deps cross_compile
		@echo "Packaging..."
		scripts/make.sh package

deps:
		scripts/deps.sh

test:
		@echo "Running tests..."
		scripts/make.sh test

clean:
		@echo "Cleaning up..."
		@rm -rf bin
		@rm -rf packages
		@rm -rf .gopath

install:
		@echo "Installing to ${INSTALLPRE}/bin"
		@scripts/make.sh install
