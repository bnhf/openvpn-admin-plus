#!/bin/bash

set -e

PKGFILE=pivpn-tap-web-ui.tar.gz

cp -f ../$PKGFILE ./

# Multi-arch the manifest way -- uncomment for each architecture, save and run script
# docker build -t bnhf/pivpn-tap-web-ui:manifest-amd64 --build-arg ARCH=amd64/ .
# docker build -t bnhf/pivpn-tap-web-ui:manifest-arm64 --build-arg ARCH=arm64/ .
# docker build -t bnhf/pivpn-tap-web-ui:manifest-armv7 --build-arg ARCH=armv7/ .

# Multi-arch the buildx way -- just use the command below, don't run the script
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -f build/Multi-arch.dockerfile -t bnhf/pivpn-tap-web-ui . --push

rm -f $PKGFILE
