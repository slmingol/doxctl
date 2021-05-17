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

netInterfacesReachableChk() {

    verbose=0; [[ $1 -eq 1 ]] && local verbose=1
    [[ $verbose -eq 1 ]] && set -x

    printf "\n\nInterfaces Reachable Checks\n===========================\n\n\n"
    scutil_nwi=$(scutil --nwi)
    net_ifs=$(echo "$scutil_nwi" | grep 'Network interfaces:' | cut -d" " -f 3-)
    inf_cnt=$(echo "$net_ifs" | wc -w | tr -d ' ')
    inf_reachable_cnt=$(echo "$scutil_nwi" \
        | grep address -B1 -A1 \
        | grep -E "flags|reach" \
        | paste - - \
        | column -t \
        | grep -v Reachable \
        | wc -l \
        | tr -d ' ')

    printf "How many network interfaces found? \t\t ---> %s <--- \t [%s]\n" "$inf_cnt" "$net_ifs"
    printf "At least 1 interface's a utun device? \t\t ---> %s <---\n" "$(echo "$net_ifs" | grep -q utun && echo yes || echo no)"
    printf "All interfaces are reachable? \t\t\t ---> %s <---\n" "$(echo "$inf_reachable_cnt" | grep -q 0 && echo yes || echo no)"

    [[ $verbose -eq 1 ]] && set +x

    echo ''
    echo ''
}

#----------------------------------------------------

vpnInterfaceRoutesChk() {

    verbose=0; [[ $1 -eq 1 ]] && local verbose=1
    [[ $verbose -eq 1 ]] && set -x

    printf "\n\nVPN Interface Route Checks\n===========================\n\n\n"
    scutil_nwi=$(scutil --nwi)
    vpn_if=$(echo "$scutil_nwi" | grep 'Network interfaces:' | grep -o utun[0-9] || echo "NIL")

    netstatOut=$(netstat -r -f inet | grep "$vpn_if")
    vpnRouteCnt=$(echo "$netstatOut" | grep "$vpn_if" -c)

    printf "At least 5 routes using interface [%s]? \t\t ---> %s <--- \t [%s]\n" \
            "$vpn_if" \
            "$([[ vpnRouteCnt -ge 5 ]] && echo yes || echo no)" \
            "$vpnRouteCnt"

    [[ $verbose -eq 1 ]] && set +x

    echo ''
    echo ''
}

#----------------------------------------------------
