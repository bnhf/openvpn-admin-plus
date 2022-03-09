# PiVPN-TAP-web-ui

## Summary
OpenVPN TAP (bridge) or TUN host-based server web administration interface in a Docker container. Intended for use with PiVPN (on amd64 versions of Debian or Ubuntu, or on ARM64/ARMv7 with Raspberry Pi OS). PiVPN should be installed first!  Recently adapted to work with host-installed TUN servers too.

Here's a post on setting PiVPN in TAP server mode. It's written for the Raspberry Pi, but the steps are the same on Debian or Ubuntu. One exception is your physical ethernet adapter, which will likely not be eth0 in the openvpn-bridge script:

https://technologydragonslayer.com/2022/01/16/installing-an-openvpn-tap-server-on-a-raspberry-pi-using-pivpn/

Also, on Debian 11 and its derivations, the bridge-conf, bridge-start and bridge-stop scripts defined here work much better than the classic openvpn-bridge script for TAP installations:

[bridge-scripts](https://gist.github.com/Belphemur/3b03eaad96172b2159fc)

No need to replace your openvpn@.service file, just insert the bridge-start and bridge-stop lines in the same relative locations in your existing file. build-client and server.conf are not needed either.

Goal: create quick to deploy and easy to use solution that makes work with small OpenVPN environments a breeze.

If you have docker and Portainer installed, you can jump directly to [installation](#Prod).

![Status page](https://user-images.githubusercontent.com/41088895/155858300-95d0b0aa-4568-42f2-9734-52a39139cf18.png)


If you have a functioning OpenVPN TAP or TUN Server on the same host as your Docker containers, you should be able
to use this fork to monitor OpenVPN connections.

Certificate generation and management is also available, and should be compatible with PiVPN. You can use either this web-ui to create client certificates, or use PiVPN from the commandline. Use PiVPN from the commandline (with elevated priveleges) to revoke certificates.

## Motivation

* to create a version of this project that will work with OpenVPN TAP and TUN servers created using PiVPN (amd64, arm64 or ARMv7 )

## Features

* status page that shows server statistics and list of connected clients
* easy creation of client certificates
* ability to download client certificates as a zip package with client configuration inside or as a single .ovpn file
* log preview
* modification of OpenVPN configuration file through web interface
* this fork is especially designed to use an external version of OpenVPN configured for TAP (bridge) -- which is probably not possible via Docker
* works with host-base PiVPN TUN servers now too!

## Screenshots

![Screenshot 2022-02-26 113330](https://user-images.githubusercontent.com/41088895/155858411-f0413188-2481-473a-891b-4e4305e3e515.png)

![screenshot-raspberrypi5_8080-2022 02 26-14_10_25](https://user-images.githubusercontent.com/41088895/155859338-b7ca2743-b702-4eff-a2d5-31144d4a1be8.png)


![Screenshot 2022-02-26 113707](https://user-images.githubusercontent.com/41088895/155858443-581b9206-327b-4df3-ac14-cd310cae768e.png)

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

Setup your Portainer Stacks page as shown on an amd64 machine running Debian or Ubuntu, inserting environment variables for creating certificates. Also, you'll need the unique ID assigned by PiVPN to the server (the name used for the server certificate and key, which is the hostname followed by a series of numbers, letters and dashes). And finally, you'll need to supply the name of the server configuration file you'd like to use (usually server.conf):

![screenshot-brix-pc2_9000-2022 03 09-13_28_03](https://user-images.githubusercontent.com/41088895/157535447-3bde8e7b-1d48-41a8-8664-d11195982283.png)

Setup your Portainer Stacks page as shown on an ARMv7 running the Raspberry Pi OS, inserting environment variables for creating certificates (grab the contents of the docker-compose.yml in the docs folder, and add a :armv7 as a tag after the container name). Also, you'll need the unique ID assigned by PiVPN to the server (the name used for the server certificate and key, which is the hostname followed by a series of numbers, letters and dashes). And finally, you'll need to supply the name of the server configuration file you'd like to use (usually server.conf). See the above amd64 screenshot for current env variable requirements.:

![screenshot-raspberrypi6_9000-2022 03 09-07_44_40](https://user-images.githubusercontent.com/41088895/157468105-f679a703-f02b-48c4-abcb-7e54b8255525.png)

Setup your Portainer Stacks page as shown on an ARM64 running the Raspberry Pi OS, inserting environment variables for creating certificates (grab the contents of the docker-compose.yml in the docs folder, and add a :arm64 as a tag after the container name). Also, you'll need the unique ID assigned by PiVPN to the server (the name used for the server certificate and key, which is the hostname followed by a series of numbers, letters and dashes). And finally, you'll need to supply the name of the server configuration file you'd like to use (usually server.conf). See the above amd64 screenshot for current env variable requirements.:

![screenshot-raspberrypi5_9000-2022 03 09-07_46_28](https://user-images.githubusercontent.com/41088895/157468989-bbb6b400-cc78-43d9-a19e-4e45d22e161a.png)

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

Portainer, Cockpit and pivpn-tap-web-ui can all be added as "vertical" tabs in organizr, for a clean single tab in your browser. All are iFrame compatible when accessed via http:// -- it'll work with https:// too, but not in iFrames.

![screenshot-nuc10-pc2-2022 03 09-14_45_09](https://user-images.githubusercontent.com/41088895/157541809-7aa4b827-9c9c-411a-9d3b-2b47fd330c31.png)

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

    go get github.com/bnhf/pivpn-tap-web-ui
    cd $GOPATH/src/github.com/bnhf/pivpn-tap-web-ui
    go mod tidy
    bee run -gendoc=true
    bee pack -exr='^vendor|^data.db|^build|^README.md|^docs'
    cd build
    ./build.sh
    
For building on ARM64 or ARMv7:

    In the dockerfile inside the build folder, comment out debian:bullseye as a source, and uncomment balenalib/raspberry-pi-debian:latest (ARMv7 only)
    In build.sh, change the docker build to <your-docker-hub-repo-here>/pivpn-tap-web-ui:arm64 (or armv7)
    It's highly recommended that you use Visual Studio Code with the "Remote - SSH" extension (in addition to the "Go" extension of course) from a more powerful machine
    

## Todo

* ARMv7 version for the Raspberry Pi -- Done!
* Update "Memory usage" on the status page to display more accurate data -- Issue reported to Cloud Foundry
* Add certificate revocation from the GUI -- currently can be done only from the commandline via PiVPN -r username (or in the Cockpit Terminal!)


## License

This project uses [MIT license](LICENSE)


## Remarks

Numerous things have been updated to bring this project forward from its 2017 roots. It's now based on Debian 11 (in the container build), and is using the latest OpenVPN and EasyRSA, thanks to PiVPN. All of the project dependencies (vendoring) have been updated to current levels in 2022.

Courtsey of @tyzbit, the ability to specify DNS servers, and additional client/server options have been added. Also @mendoza-conicet contributed code for being able to download a single .ovpn file. Many issues have been addressed related to adapting this package for use with a host-based server, and related to all of the latest versions of the dependencies.

And, of course, many thanks to @adamwalach for his excellent original work to create this project!


### Template
AdminLTE - dashboard & control panel theme. Built on top of Bootstrap 3.
