#!/bin/bash -e

# PiVPN will have already setup the CA and server.crt/.key
# CA_NAME=LocalCA
# SERVER_NAME=server
# EASY_RSA=/etc/openvpn/easy-rsa

VARS=/etc/openvpn/easy-rsa/vars
# PiVPN stores its keys under easy-rsa/pki
# mkdir -p /etc/openvpn/keys
# touch /etc/openvpn/keys/index.txt
# echo 01 > /etc/openvpn/keys/serial
cp -f /opt/scripts/vars.template /etc/openvpn/easy-rsa/vars

# Append the env variables passed by Docker to the vars file
echo -e "\n"                                       >> $VARS
echo "set_var EASYRSA_REQ_COUNTRY   \"$COUNTRY\""  >> $VARS
echo "set_var EASYRSA_REQ_PROVINCE  \"$PROVINCE\"" >> $VARS
echo "set_var EASYRSA_REQ_CITY      \"$CITY\""     >> $VARS
echo "set_var EASYRSA_REQ_ORG       \"$ORG\""      >> $VARS
echo "set_var EASYRSA_REQ_EMAIL     \"$EMAIL\""    >> $VARS
echo "set_var EASYRSA_REQ_OU        \"$OU\""       >> $VARS

# Append name=server to the end of the first line of index.txt

sed -i ' 1 s/.*/&\/name=server/' /etc/openvpn/easy-rsa/pki/index.txt

# Determine commonname for server and add to environment

# $EASY_RSA/clean-all
# source /etc/openvpn/keys/vars
# export KEY_NAME=$CA_NAME
# echo "Generating CA cert"
# $EASY_RSA/build-ca
# export EASY_RSA="${EASY_RSA:-.}"

# $EASY_RSA/easyrsa --batch build-ca nopass $*

# export KEY_NAME=$SERVER_NAME

# echo "Generating server cert"
# $EASY_RSA/build-key-server $SERVER_NAME
# $EASY_RSA/easyrsa --batch build-server-full $SERVER_NAME nopass
