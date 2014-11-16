#!/bin/bash

case $1 in
"test")
  pushd fastfood &> /dev/null
  go test
  popd &> /dev/null
  pushd fastfood/cookbook &> /dev/null
  go test
  popd &> /dev/null
  ;;
"build")
  mkdir -p ${GOPATH}/src/github.com/jarosser06
  go build -o bin/fastfood ./fastfood.go
  ;;
esac
