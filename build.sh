#!/bin/bash
set -e

BUILD_VERSION="0.01"

echo "=== Test === "
go test -cover ./...
go build -o halfpipe -ldflags "-X main.version=${BUILD_VERSION}" cmd/halfpipe.go

echo; echo "=== Smoke Test ==="
set +e
RESULT=$(./halfpipe --version 2>&1)

echo "${RESULT}"

if [ "${RESULT}" == "halfpipe ${BUILD_VERSION}" ]; then
   echo ok
else
    echo Failed:
    exit 1
fi
