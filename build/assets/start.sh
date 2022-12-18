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
cd /opt/openvpn-gui-tap
echo "Working directory set to" $PWD

if [ ! -z $ENABLEHTTPS ]; then
  sed -i '/EnableHTTPS=/s/.*/EnableHTTPS='"$ENABLEHTTPS"'/' conf/app.conf
  echo "HTTPS enabled set to \"$ENABLEHTTPS\""
fi

if [ ! -z $HTTPSPORT ]; then
  sed -i '/HTTPSPort=/s/.*/HTTPSPort='"$HTTPSPORT"'/' conf/app.conf
  echo "HTTPS port set to: \"$HTTPSPORT\""
fi

if [ ! -z $HTTPSCERT ]; then
  sed -i '/HTTPSCertFile=/s/.*/HTTPSCertFile='"$HTTPSCERT"'/' conf/app.conf
  echo "HTTPS Certificate path set to: \"$HTTPSCERT\""
else
  sed -i '/HTTPSCertFile=/s/.*/HTTPSCertFile=\/etc\/openvpn\/easy-rsa\/pki\/issued\/'"$PIVPN_SERVER"'.crt/' conf/app.conf
  echo "HTTPS Certificate path set to default: \"$HTTPSCERT\""
fi

if [ ! -z $HTTPSKEY ]; then
  sed -i '/HTTPSKeyFile=/s/.*/HTTPSKeyFile='"$HTTPSKEY"'/' conf/app.conf
  echo "HTTPS key path set to: \"$HTTPSKEY\""
else
  sed -i '/HTTPSKeyFile=/s/.*/HTTPSKeyFile=\/etc\/openvpn\/easy-rsa\/pki\/private\/'"$PIVPN_SERVER"'.key/' conf/app.conf
  echo "HTTPS key path set to default: \"$HTTPSKEY\""
fi

mkdir -p db
echo "Starting!"
./pivpn-tap-web-ui


