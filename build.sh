#!/bin/sh
GOBUILDFLAGS="-v -trimpath -ldflags -s"

mkdir -p out
go build $GOBUILDFLAGS -o out/lets-hit-the-gym .
