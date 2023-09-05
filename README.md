# OpenVPN-Admin-Plus

## Summary
OpenVPN TAP/bridge or TUN (host-based) server web administration interface in a Docker container. Intended for use with PiVPN (on amd64/arm64/armv7 versions of Debian or Ubuntu, or on ARM64/ARMv7 with Raspberry Pi OS). PiVPN should be installed first and OpenVPN 2.5.x or above required!

Here's a post on setting PiVPN in TAP server mode. It's written for the Raspberry Pi, but the steps are the same on Debian or Ubuntu. One exception is your physical ethernet adapter, which will likely not be eth0 in the openvpn-bridge script:

https://technologydragonslayer.com/2022/01/16/installing-an-openvpn-tap-server-on-a-raspberry-pi-using-pivpn/

IMPORTANT: On Debian 11 and its derivations, the bridge-conf, bridge-start and bridge-stop scripts defined here work much better than the classic openvpn-bridge script for TAP installations:

[bridge-scripts](https://gist.github.com/Belphemur/3b03eaad96172b2159fc)

No need to replace your openvpn@.service file, just insert the bridge-start and bridge-stop lines in the same relative locations in your existing file. build-client and server.conf are not needed either.

Goal: create quick to deploy and easy to use solution that makes work with small OpenVPN environments a breeze.

If you have docker and Portainer installed, you can jump directly to [installation](#Prod).

![Status page](https://user-images.githubusercontent.com/41088895/211678414-57d18b6e-0ba0-4f30-b7d5-deebdc592252.png)

If you have a functioning OpenVPN TAP or TUN Server on the same host as your Docker containers, you should be able to use this fork to monitor OpenVPN connections.

Certificate generation and management is also available, and should be compatible with PiVPN. You can use either this web-ui to create client certificates, or use PiVPN from the commandline. Use PiVPN from the commandline (with elevated privileges) to revoke certificates.

## Motivation

* to create a version of this project that will work with OpenVPN TAP and TUN servers created using PiVPN (amd64, arm64 or ARMv7 )

## Features

* status page that shows server statistics and list of connected clients
* easy creation of client certificates
* ability to download client certificates as a zip package with client configuration inside or as a single .ovpn file
* log preview
* modification of OpenVPN configuration file through web interface
* this fork is especially designed to use an external version of OpenVPN configured for TAP (bridge) -- which is probably not possible via Docker
* works with host-based PiVPN TUN servers now too!

## Screenshots

![screenshot-nuc10-pc2 localdomain-2023 01 10-14_51_43](https://user-images.githubusercontent.com/41088895/211678711-b41de1c8-590f-45c0-9048-091ed1265211.png)

![screenshot-nuc10-pc2 localdomain_8080-2023 01 10-15_11_54](https://user-images.githubusercontent.com/41088895/211678760-ed8a8425-21db-4f28-b6f0-51078450d2d8.png)

![screenshot-nuc10-pc2 localdomain-2023 01 10-14_45_55](https://user-images.githubusercontent.com/41088895/211678922-c107f5d8-b642-40f3-8f5a-ae3a2cd6d662.png)

![Screenshot 2022-02-26 113822](https://user-images.githubusercontent.com/41088895/155858448-cced00d9-b931-4e85-a77f-f0f220ac0afc.png)

[Screenshots](docs/screenshots.md)

## Usage

After startup web service is visible on port 8080. To login use the following default credentials:

* username: admin
* password: b3secure

Please change password to your own immediately!

### Prod

Requirements:
* Docker, Portainer, PiVPN, Debian or Ubuntu
* on firewall open ports: 8080/tcp

Setup your Portainer Stacks page as shown on an amd64/arm64/armv7 machine running Debian or its derivations, inserting environment variables for creating certificates. Also, you'll need the unique ID assigned by PiVPN to the server (the name used for the server certificate and key, which is the hostname followed by a series of numbers, letters and dashes). And finally, you'll need to supply the name of the server configuration file you'd like to use (usually server.conf):

![screencapture-brix-pc2-9000-2023-02-04-14_15_54](https://user-images.githubusercontent.com/41088895/216790336-32dbaebe-7f6d-4c82-8b13-db55fc59c7d6.png)

Copy and paste this into the Portainer Environment variables section in Advanced Mode, then switch to Simple mode to input your values. Replace the ${variable} with your actual values, and remove the comment portions that begin with #. Leave USERNAME and PASSWORD default values, as you'll change those after you login for the first time by clicking on Admistrator - Profile:

```yml
 OPENVPN_ADMIN_USERNAME=admin # Leave this default as-is and update on first-run
 OPENVPN_ADMIN_PASSWORD=b3secure # Leave this default as-is and update on first-run
 PIVPN_CONF=server.conf # Filename for your OpenVPN server configuration
 COUNTRY=${COUNTRY} # Two character country code. Country, Province, City, Org and OU are used by EasyRSA
 PROVINCE=${PROVINCE} # Province or State
 CITY=${CITY}
 ORG=${ORG} # Organization - Use your company name here or make one up
 EMAIL=${EMAIL}
 OU=${OU} # Organizational Unit - Use your company department name or make one up
 TZ=${TZ} # Your OpenVPN server's timezone
 ENABLEHTTPS=${ENABLEHTTPS} # Enable HTTPS protocol
 HTTPSPORT=${HTTPSPORT} # Specify the HTTPS port number
 HTTPSCERT=${HTTPSCERT} # The path to the SSL-certificate (for example: /etc/openvpn/server/ssl/openvpn-server.crt)
 HTTPSKEY=${HTTPSKEY} # The path to the private key (for example: /etc/openvpn/server/ssl/openvpn-server.key)
```

This fork uses a single docker container with the OpenVPNAdmin web application. Through a docker volume it creates following directory structure for the database, but otherwise links to /etc/openvpn in the host. The intention is for PiVPN to be able to operate as usual, with PiVPN commanline options still available:

    .
    ├── docker-compose.yml
    └── openvpn-data
         └── db
            └── data.db

### User

Requirements:
* [docker](https://docs.docker.com/engine/install/debian/#install-using-the-convenience-script)

Optional, but highly recommended:
* [Portainer](https://docs.portainer.io/v/ce-2.9/start/install/server/docker/linux)
* [cockpit-project](https://cockpit-project.org)
* [cockpit-navigator plugin](https://cockpit-project.org/applications)
* [organizr-Docker](https://hub.docker.com/r/organizr/organizr)
* [watchtower](https://hub.docker.com/r/containrrr/watchtower)

Portainer, Cockpit and pivpn-tap-web-ui can all be added as "vertical" tabs in organizr, for a clean single tab in your browser. All are iFrame compatible when accessed via http:// -- it'll work with https:// too, but not in the iFrame (new tab will open). watchtower can be setup to run once and then stop, updating all of your containers (including Portainer and watchtower itself).

![screenshot-brix-pc2-2022 03 09-14_50_38](https://user-images.githubusercontent.com/41088895/157542989-5f13f2e5-8b69-4958-a3dc-95270485efc0.png)

### Dev

Requirements:
* [golang environments](https://www.digitalocean.com/community/tutorial_series/how-to-code-in-go)
* [beego](https://beego.vip/)
* [bee](https://github.com/beego/bee)
* [docker](https://docs.docker.com/engine/install/debian/#install-using-the-convenience-script)

Optional, but recommended:

* [Portainer](https://docs.portainer.io/v/ce-2.9/start/install/server/docker/linux)
* [GitHub Desktop for Linux](https://gist.github.com/berkorbay/6feda478a00b0432d13f1fc0a50467f1)
* [Visual Studio Code](https://code.visualstudio.com/download)

Execute commands:

    go get github.com/bnhf/openvpn-admin-plus
    cd $GOPATH/src/github.com/bnhf/openvpn-admin-plus
    go mod tidy
    bee run -gendoc=true
    bee pack -exr='^vendor|^data.db|^build|^README.md|^docs'
    cd build
    ./build.sh
    
For building on ARM64 or ARMv7:

    In the dockerfile inside the build folder, comment out debian:bullseye as a source,
    and uncomment balenalib/raspberry-pi-debian:latest (ARMv7 only)
    In build.sh, change the docker build to <your-docker-hub-repo-here>/pivpn-tap-web-ui:arm64 (or armv7)
    It's highly recommended that you use Visual Studio Code with the "Remote - SSH" extension
    (in addition to the "Go" extension of course) from a more powerful machine
    

## Todo

* Test with TUN and TAP running concurrently

## License

This project uses [MIT license](LICENSE)


## Remarks

Numerous things have been updated to bring this project forward from its 2017 roots. It's now based on Debian 11 (in the container build), and is using the latest OpenVPN and EasyRSA, thanks to PiVPN. All of the project dependencies (vendoring) have been updated to current levels in 2022.

Courtsey of @tyzbit, the ability to specify DNS servers, and additional client/server options have been added. Also @mendoza-conicet contributed code for being able to download a single .ovpn file. Many issues have been addressed related to adapting this package for use with a host-based server, and related to all of the latest versions of the dependencies.

And, of course, many thanks to @adamwalach for his excellent original work to create this project!


### Template
AdminLTE - dashboard & control panel theme. Built on top of Bootstrap 3.
