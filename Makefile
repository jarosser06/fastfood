export GOPATH := ${PWD}/.gopath
export RICEBIN := ${GOPATH}/bin/rice

all: deps build

build:
		@echo "Building binary..."
		scripts/make.sh build

deps:
		scripts/deps.sh

rice:
		@echo "Compiling templates..."
		go get github.com/GeertJohan/go.rice/rice
		go install github.com/GeertJohan/go.rice/rice
		${RICEBIN} -i ./provider/cookbook clean
		${RICEBIN} -i ./provider/cookbook embed-go
		${RICEBIN} -i ./provider/application clean
		${RICEBIN} -i ./provider/application embed-go

test:
		@echo "Running tests..."
		scripts/make.sh test

clean:
		@echo "Cleaning up..."
		rm -rf bin
		rm -rf .gopath

install:
		@echo "Installing to user's bin..."
		scripts/make.sh install
