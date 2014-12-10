#!/bin/bash

case $1 in
"test")
  pkg_dirs="cmd"

  for dir in $pkg_dirs
  do
    pushd ${dir} &> /dev/null
    go test
    popd &> /dev/null
  done

  go test
  ;;
"build")
  mkdir -p ${GOPATH}/src/github.com/jarosser06
  go build -o bin/fastfood ./cmd/main/main.go
  ;;
"install")
  cp bin/fastfood /usr/local/bin
  cp doc/manpage /usr/local/share/man/man1/fastfood.1
  ;;
esac
