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
  sed -i 's/EnableHTTPS=false/EnableHTTPS=true/g' conf/app.conf
  echo "HTTPS enabled"
fi

if [ ! -z $HTTPSPORT ]; then
  sed -i 's/HTTPSPort=8443/HTTPSPort='"$HTTPSPort"'/g' conf/app.conf
  echo "HTTPS port set to: \"$HTTPSPORT\""
fi

if [ ! -z $HTTPSCERT ]; then
  sed -i 's/HTTPSCertFile=/HTTPSCertFile='"$HTTPSCERT"'/g' conf/app.conf
  echo "HTTPS Certificate path set to: \"$HTTPSCERT\""
else
  sed -i 's/HTTPSCertFile=/HTTPSCertFile=\/etc\/openvpn\/easy-rsa\/pki\/issued\/'"$PIVPNSERVER"'.crt /g' conf/app.conf
fi

if [ ! -z $HTTPSKEY ]; then
  sed -i 's/HTTPSKeyFile=/HTTPSKeyFile='"$HTTPSKEY"'/g' conf/app.conf
  echo "HTTPS key path set to: \"$HTTPSKEY\""
else
  sed -i 's/HTTPSKeyFile=/HTTPSKeyFile=\/etc\/openvpn\/easy-rsa\/pki\/private\/'"$PIVPNSERVER"'.key /g' conf/app.conf
fi

mkdir -p db
echo "Starting!"
./pivpn-tap-web-ui


