#!/bin/bash

pkg_dirs="util cookbook config helpers application"

case $1 in
"test")
  for dir in $pkg_dirs
  do
    pushd pkg/${dir} &> /dev/null
    go test
    popd &> /dev/null
  done
  ;;
"build")
  mkdir -p ${GOPATH}/src/github.com/jarosser06
  go build -o bin/fastfood ./fastfood.go
  ;;
"install")
  cp bin/fastfood /usr/local/bin
  ;;
esac
