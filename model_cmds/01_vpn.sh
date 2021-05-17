#!/usr/bin/env bash

pasteCmd="gpaste"

### MacOS

###### OFF VPN
###
### $ scutil --nwi
### Network information
###
### IPv4 network interface information
###      en0 : flags      : 0x5 (IPv4,DNS)
###            address    : 192.168.13.10
###            reach      : 0x00000002 (Reachable)
###
###    REACH : flags 0x00000002 (Reachable)
###
### IPv6 network interface information
###    No IPv6 states found
###
###
###    REACH : flags 0x00000000 (Not Reachable)
###
### Network interfaces: en0


##### ON VPN
### $ scutil --nwi
### Network information
###
### IPv4 network interface information
###      en0 : flags      : 0x5 (IPv4,DNS)
###            address    : 192.168.13.10
###            reach      : 0x00000002 (Reachable)
###    utun2 : flags      : 0x7 (IPv4,IPv6,DNS)
###            address    : 172.31.250.118
###            reach      : 0x00000002 (Reachable)
###
###    REACH : flags 0x00000002 (Reachable)
###
### IPv6 network interface information
###    utun2 : flags      : 0x7 (IPv4,IPv6,DNS)
###            address    : 2001:db8::b
###            reach      : 0x00000002 (Reachable)
###
###    REACH : flags 0x00000002 (Reachable)
###
### Network interfaces: en0 utun2

#----------------------------------------------------



dnsResolverChk() {

 verbose=0; [[ $1 -eq 1 ]] && local verbose=1
 [[ $verbose -eq 1 ]] && set -x

 [[ $verbose -eq 1 ]] && set +x
 echo ''
 echo ''
}

#----------------------------------------------------
