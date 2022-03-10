docker manifest create \
bnhf/pivpn-tap-web-ui:manifest-latest \
--amend bnhf/pivpn-tap-web-ui:manifest-amd64 \
--amend bnhf/pivpn-tap-web-ui:manifest-armv7 \
--amend bnhf/pivpn-tap-web-ui:manifest-arm64
