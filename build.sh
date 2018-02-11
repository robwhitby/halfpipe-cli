#!/bin/bash
set -x -e

go test -cover ./...
go install
