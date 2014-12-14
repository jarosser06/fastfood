export GOPATH := ${PWD}/.gopath

all: deps build

build:
		@echo "Building binaries..."
		go get github.com/mitchellh/gox
		${GOPATH}/bin/gox -os="linux darwin windows" -output="bin/fastfood_{{.OS}}_{{.Arch}}" ./cmd/main/

deps:
		scripts/deps.sh

test:
		@echo "Running tests..."
		scripts/make.sh test

clean:
		@echo "Cleaning up..."
		rm -rf bin
		rm -rf .gopath

install:
		@echo "Installing to user's bin..."
		@scripts/make.sh install
