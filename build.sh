#!/bin/sh
GOBUILDFLAGS="-v -trimpath -ldflags -s"
TARGET_DIR="out/$(go env GOOS)/$(go env GOARCH)"

mkdir -p $TARGET_DIR
go build $GOBUILDFLAGS -o $TARGET_DIR/lets-hit-the-gym .
