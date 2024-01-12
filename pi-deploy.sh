#!/bin/sh
GOOS=linux GOARCH=arm64 ./build.sh
echo "put out/linux/arm64/lets-hit-the-gym" | sftp lets-hit-the-gym@pi.quantaly.net
