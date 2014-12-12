export GOPATH := ${PWD}/.gopath

all: deps build

build:
		@echo "Building binary..."
		scripts/make.sh build

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
