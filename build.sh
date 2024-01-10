#!/bin/sh
GOBUILDFLAGS="-v -trimpath -ldflags -s"

mkdir -p out
go build $GOBUILDFLAGS -o out/remind ./bin/remind 
go build $GOBUILDFLAGS -o out/schedule ./bin/schedule
