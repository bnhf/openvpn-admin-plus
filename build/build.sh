#!/bin/bash

set -e

PKGFILE=pivpn-tap-web-ui.tar.gz

cp -f ../$PKGFILE ./

docker build -t bnhf/pivpn-tap-web-ui:manifest-amd64 --build-arg ARCH=amd64/ .
# docker build -t bnhf/pivpn-tap-web-ui:manifest-arm64 --build-arg ARCH=arm64/ .
# docker build -t bnhf/pivpn-tap-web-ui:manifest-armv7 --build-arg ARCH=armv7/ .

rm -f $PKGFILE
