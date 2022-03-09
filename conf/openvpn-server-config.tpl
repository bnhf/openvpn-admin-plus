management {{ .Management }}
dev {{ .Dev }}
proto {{ .Proto }}
port {{ .Port }}

ca {{ .Ca }}
cert {{ .Cert }}
key {{ .Key }}
dh {{ .Dh }}
ecdh-curve prime256v1

topology subnet
{{ .Server }}
ifconfig-pool-persist {{ .IfconfigPoolPersist }}
push "dhcp-option DNS {{ .DNSServerOne }}"
push "dhcp-option DNS {{ .DNSServerTwo }}"

keepalive {{ .Keepalive }}
remote-cert-tls client
tls-version-min 1.2
tls-crypt {{ .CCEncryption }}
cipher {{ .Cipher }}
auth {{ .Auth }}

persist-key
persist-tun
crl-verify /etc/openvpn/crl.pem

# status /etc/openvpn/openvpn-status.log 20
# status-version 3
# syslog
log /etc/openvpn/openvpn.log
verb 3
mute 10

{{ .ExtraServerOptions }}