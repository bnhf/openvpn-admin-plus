#!/bin/bash

set -e

PKGFILE=openvpn-tap-external-web-ui.tar.gz

cp -f ../$PKGFILE ./

docker build -t bnhf/openvpn-tap-external-web-ui .

rm -f $PKGFILE
