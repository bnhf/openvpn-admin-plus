FROM bnhf/go-beego-bee-git
WORKDIR /go/src/github.com/bnhf

# Uncomment for a multi-arch buildx of the main branch
RUN git clone https://github.com/bnhf/pivpn-tap-web-ui
# Uncomment for a multi-arch buildx of the develop branch
# RUN git clone -b develop --single-branch https://github.com/bnhf/pivpn-tap-web-ui
WORKDIR /go/src/github.com/bnhf/pivpn-tap-web-ui
RUN go mod tidy && \
    bee pack -exr='^vendor|^data.db|^build|^README.md|^docs'

FROM debian:bullseye
WORKDIR /opt
EXPOSE 8080

RUN apt-get update && apt-get install -y easy-rsa && \
    chmod 755 /usr/share/easy-rsa/*
COPY --from=0  /go/src/github.com/bnhf/pivpn-tap-web-ui/build/assets/start.sh /opt/start.sh
COPY --from=0  /go/src/github.com/bnhf/pivpn-tap-web-ui/build/assets/generate_ca_and_server_certs.sh /opt/scripts/generate_ca_and_server_certs.sh
COPY --from=0  /go/src/github.com/bnhf/pivpn-tap-web-ui/build/assets/vars.template /opt/scripts/

COPY --from=0 /go/src/github.com/bnhf/pivpn-tap-web-ui/pivpn-tap-web-ui.tar.gz /opt/openvpn-gui-tap/
RUN tar -zxf /opt/openvpn-gui-tap/pivpn-tap-web-ui.tar.gz --directory /opt/openvpn-gui-tap/
RUN rm -f /opt/openvpn-gui-tap/db/data.db /opt/openvpn-gui-tap/pivpn-tap-web-ui.tar.gz
COPY --from=0 /go/src/github.com/bnhf/pivpn-tap-web-ui/build/assets/app.conf /opt/openvpn-gui-tap/conf/app.conf

CMD /opt/start.sh
