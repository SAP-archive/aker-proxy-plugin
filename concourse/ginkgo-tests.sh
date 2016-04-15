#!/bin/bash

set -e

mkdir -p $GOPATH/src

echo "Relocating vendored packages..."
mv aker-proxy/vendor/* $GOPATH/src

echo "Relocating project..."
mkdir -p $GOPATH/src/github.infra.hana.ondemand.com/I061150
cp -r aker-proxy $GOPATH/src/github.infra.hana.ondemand.com/I061150
cd $GOPATH/src/github.infra.hana.ondemand.com/I061150/aker-proxy

echo "Fetching test tools..."
go get github.com/onsi/ginkgo/ginkgo

echo "Running tests..."
ginkgo -r
