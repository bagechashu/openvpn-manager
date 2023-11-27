# OpenVPN manager by web console

### How to use
1. Exec script/openvpn-install.sh to install openvpn. (the script is maintenance by [Nyr](https://github.com/Nyr/openvpn-install) )
2. Exec `go run cmd/openvpn-manager/main.go -c openvpn-manager.ini.` (Of course, we can compile to binary first.)
3. visit http://127.0.0.1:8080

### Main functional
- Add user
- Revoke user
- Show user
- Download user config

# Thanks
> https://github.com/Nyr/openvpn-install

### openvpn-install
OpenVPN [road warrior](http://en.wikipedia.org/wiki/Road_warrior_%28computing%29) installer for Ubuntu, Debian, AlmaLinux, Rocky Linux, CentOS and Fedora.

This script will let you set up your own VPN server in no more than a minute, even if you haven't used OpenVPN before. It has been designed to be as unobtrusive and universal as possible.

#### Installation
Run the script and follow the assistant:

`wget https://raw.githubusercontent.com/Nyr/openvpn-install/master/openvpn-install.sh -O openvpn-install.sh && bash openvpn-install.sh`

Once it ends, you can run it again to add more users, remove some of them or even completely uninstall OpenVPN.
