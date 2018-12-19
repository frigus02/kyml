#!/bin/bash

set -eo pipefail

if [ "$TRAVIS_TAG" != "" ]; then
	if [[ $TRAVIS_TAG != v* ]]; then
		echo >&2 "TRAVIS_TAG does not start with \"v\""
		exit 1
	fi
	version="${TRAVIS_TAG:1}"
else
	version=${1:?"Version (arg 1) missing"}
fi

echo "Building version $version..."

ldflags="-X github.com/frigus02/kyml/pkg/commands.version=$version"

GOARCH=amd64 GOOS=darwin go build -o "bin/kyml_${version}_darwin_amd64" -ldflags "$ldflags"
GOARCH=amd64 GOOS=linux go build -o "bin/kyml_${version}_linux_amd64" -ldflags "$ldflags"
GOARCH=amd64 GOOS=windows go build -o "bin/kyml_${version}_windows_amd64.exe" -ldflags "$ldflags"
