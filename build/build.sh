#!/bin/bash

set -e

PKGFILE=pivpn-tap-web-ui.tar.gz

cp -f ../$PKGFILE ./

docker build -t bnhf/pivpn-tap-web-ui .

rm -f $PKGFILE
