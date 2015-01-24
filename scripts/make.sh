#!/bin/bash

case $1 in
"test")
  pkg_dirs="common/fileutil common/stringutil common/maputil framework framework/chef"

  go test

  for dir in $pkg_dirs
  do
    pushd ${dir} &> /dev/null
    go test
    popd &> /dev/null
  done
  ;;
"build")
  mkdir -p ${GOPATH}/src/github.com/jarosser06
  go build -o bin/fastfood ./cmd/main/main.go
  ;;
"package")
  if [ -d packages ]; then
    rm -rf packages/*
  else
    mkdir packages
  fi

  for dir in $(ls bin/)
  do
    cp doc/manpage bin/${dir}
    cp CHANGELOG.md bin/${dir}
    pushd bin &> /dev/null
    zip ../packages/fastfood-${dir}-${FFVERSION}.zip ${dir}/*
    tar -zcf ../packages/fastfood-${dir}-${FFVERSION}.tar.gz ${dir}/*
    popd &> /dev/null
  done
  ;;
"install")
  cp doc/manpage ${INSTALLPRE}/share/man/man1/fastfood.1
  cp bin/fastfood ${INSTALLPRE}/bin/

  ;;
esac
