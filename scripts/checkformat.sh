#!/bin/bash

set -eo pipefail

diff="$(gofmt -d -s .)"
if [ -n "$diff" ]; then
	echo "Format errors:" >&2
	echo "$diff" >&2
	echo "Please run \"gofmt -s -w .\" to correct them." >&2
	exit 1
fi
