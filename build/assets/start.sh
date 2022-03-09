#!/bin/bash

clear
set -e
OVDIR=/etc/openvpn
echo "OpenVPN directory set to" $OVDIR

cd /opt/
echo "Working Directory set to" $PWD

if [ ! -f $OVDIR/.provisioned ]; then
  echo "Preparing vars"
  mkdir -p $OVDIR
  ./scripts/generate_ca_and_server_certs.sh
#  openssl dhparam -dsaparam -out $OVDIR/dh2048.pem 2048
  touch $OVDIR/.provisioned
fi

export PIVPN_SERVER=$(awk -F= '/server/ {print $2}' \
  /etc/openvpn/easy-rsa/pki/index.txt \
  | awk -F/ '{print $1}')

echo "PiVPN Server set to" $PIVPN_SERVER
cd /opt/openvpn-gui
echo "Working directory set to" $PWD
mkdir -p db
echo "Starting!"
./openvpn-tap-external-web-ui


