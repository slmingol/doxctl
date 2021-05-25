#!/usr/bin/env bash

pasteCmd="gpaste"

### MacOS

### $ printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil
### <dictionary> {
###   DomainName : bandwidth.local
###   SearchDomains : <array> {
###     0 : bandwidth.local
###   }
###   SearchOrder : 1
###   ServerAddresses : <array> {
###     0 : 10.5.0.18
###     1 : 10.5.0.19
###     2 : 192.168.7.85
###   }
###   SupplementalMatchDomains : <array> {
###     0 :
###     1 : bandwidth.local
###   }
### }

#----------------------------------------------------

dnsResolverChk() {

    verbose=0; [[ $1 -eq 1 ]] && local verbose=1
    [[ $verbose -eq 1 ]] && set -x

    printf "\n\nDNS Resolver Checks\n===================\n\n\n"

    vpn_resolvers=$(printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil)
    (
        echo "$vpn_resolvers" | grep -q 'DomainName.*bandwidth.local' && echo "DomainName set" || echo "DomainName unset"
        echo "$vpn_resolvers" | grep -A1 'SearchDomains' | grep -qE '[0-1].*bandwidth' && echo "SearchDomains set" || echo "SearchDomains unset"
        echo "$vpn_resolvers" | grep -A3 'ServerAddresses' | grep -qE '[0-1].*10.5' && echo "ServerAddresses set" || echo "ServerAddresses unset"
    ) | column -t
    [[ $verbose -eq 1 ]] && set +x
    echo ''
    echo ''
}

#----------------------------------------------------

dnsResolverPingChk() {

    verbose=0; [[ $1 -eq 1 ]] && local verbose=1
    [[ $verbose -eq 1 ]] && set -x

    printf "\n\nPing Resolver Checks\n====================\n\n\n"

    vpn_resolvers=$(printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil)
    resolver_ips=$(echo "$vpn_resolvers" | grep -A3 'ServerAddresses' | grep -E '[0-1].*10.5' | cut -d':' -f2 | $pasteCmd -s -d '')

    printf "How many resolvers found? \t\t ---> %s <---\n\n" $(echo "$resolver_ips" | wc -w)

    for resolver in $resolver_ips; do
        output=$(ping -c5 ${resolver} -W 200 -t 30 -q 2>&1)
        reachable=""
        if [ $? -eq 0 ]; then
            reachable="yes"
        else
            reachable="no"
        fi

        printf "Was resolver %s pingable? \t ---> %s <---\n" "$resolver" "$reachable"

        if [ "$reachable" == "yes" ]; then
            ncTcpOutput=$(nc -z -v -w5 $resolver 53 2>&1) && tcpReachable="yes" || tcpReachable="no"
            printf "Can we reach port 53 (DNS) via TCP? \t ---> %s <--- \t\t [ %s ]\n" "$tcpReachable" "$ncTcpOutput"
            ncUdpOutput=$(nc -z -u -v -w5 $resolver 53 2>&1) && udpReachable="yes" || udpReachable="no"
            printf "Can we reach port 53 (DNS) via UDP? \t ---> %s <--- \t\t [ %s ]\n" "$udpReachable" "$ncUdpOutput"
        fi

        echo ''
        echo ''
    done

    [[ $verbose -eq 1 ]] && set +x
    echo ''
}

#----------------------------------------------------

dnsResolverDigChk() {

    verbose=0; [[ $1 -eq 1 ]] && local verbose=1
    [[ $verbose -eq 1 ]] && set -x

    printf "\n\nDig Resolver Checks\n===================\n\n\n"

    vpn_resolvers=$(printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil)
    resolver_ips=$(echo "$vpn_resolvers" | grep -A3 'ServerAddresses' | grep -E '[0-1].*10.5' | cut -d':' -f2 | $pasteCmd -s -d '')
    sites="lab1 rdu1 atl1 dfw1 lax2 jfk1"

    # should find 6 * 2 * <x resolvers> = 12 * x ||| (6 sites w/ 2 svrs in each and typically 2 resolvers so 24!)

    svrCnt=0
    for site in $sites; do
        for resolver in $resolver_ips; do

            hitCnt1=$(dig @${resolver} idm-01a.${site}.bandwidthclec.local +short | head -1 | wc -l | tr -d ' ')
            hitCnt2=$(dig @${resolver} idm-01b.${site}.bandwidthclec.local +short | head -1 | wc -l | tr -d ' ')
            printf "\tresolver: %s | site: %s | idm-01a | cnt: %s | idm-01b |  cnt: %s | \n" \
                "$resolver" "$site" "$hitCnt1" "$hitCnt2"
            [ $hitCnt1 -eq 1 ] && ((svrCnt++))
            [ $hitCnt2 -eq 1 ] && ((svrCnt++))

        done
    done

    echo ''
    echo ''

    [ $svrCnt -eq 24 ] && idmSvrsResolvable="yes" || idmSvrsResolvable="no"
    printf "Check if we can resolve all 24 IDM server names (idm-01[ab].*)? \t ---> %s <--- \t\t [ Actual: %s ]\n" "$idmSvrsResolvable" "$svrCnt"

    [[ $verbose -eq 1 ]] && set +x

    echo ''
    echo ''
}

#dnsResolverChk
#dnsResolverPingChk
#dnsResolverDigChk


###  $ scutil --dns
###  DNS configuration
###
###  resolver #1
###    search domain[0] : bandwidth.local
###    search domain[1] : bub.lan
###    nameserver[0] : 10.5.0.18
###    nameserver[1] : 10.5.0.19
###    nameserver[2] : 192.168.7.85
###    flags    : Supplemental, Request A records, Request AAAA records
###    reach    : 0x00000002 (Reachable)
###    order    : 1
###
###  resolver #2
###    nameserver[0] : 192.168.7.85
###    if_index : 6 (en0)
###    flags    : Request A records, Request AAAA records
###    reach    : 0x00000002 (Reachable)
###    order    : 200000
###
###  resolver #3
###    domain   : bandwidth.local
###    nameserver[0] : 10.5.0.18
###    nameserver[1] : 10.5.0.19
###    nameserver[2] : 192.168.7.85
###    flags    : Supplemental, Request A records, Request AAAA records
###    reach    : 0x00000002 (Reachable)
###    order    : 1
###
###  resolver #4
###    domain   : local
###    options  : mdns
###    timeout  : 5
###    flags    : Request A records, Request AAAA records
###    reach    : 0x00000000 (Not Reachable)
###    order    : 300000
###
###  resolver #5
###    domain   : 254.169.in-addr.arpa
###    options  : mdns
###    timeout  : 5
###    flags    : Request A records, Request AAAA records
###    reach    : 0x00000000 (Not Reachable)
###    order    : 300200
###
###  resolver #6
###    domain   : 8.e.f.ip6.arpa
###    options  : mdns
###    timeout  : 5
###    flags    : Request A records, Request AAAA records
###    reach    : 0x00000000 (Not Reachable)
###    order    : 300400
###
###  resolver #7
###    domain   : 9.e.f.ip6.arpa
###    options  : mdns
###    timeout  : 5
###    flags    : Request A records, Request AAAA records
###    reach    : 0x00000000 (Not Reachable)
###    order    : 300600
###
###  resolver #8
###    domain   : a.e.f.ip6.arpa
###    options  : mdns
###    timeout  : 5
###    flags    : Request A records, Request AAAA records
###    reach    : 0x00000000 (Not Reachable)
###    order    : 300800
###
###  resolver #9
###    domain   : b.e.f.ip6.arpa
###    options  : mdns
###    timeout  : 5
###    flags    : Request A records, Request AAAA records
###    reach    : 0x00000000 (Not Reachable)
###    order    : 301000
###
###  DNS configuration (for scoped queries)
###
###  resolver #1
###    search domain[0] : bub.lan
###    nameserver[0] : 192.168.7.85
###    if_index : 6 (en0)
###    flags    : Scoped, Request A records
###    reach    : 0x00000002 (Reachable)
###
###  resolver #2
###    search domain[0] : bandwidth.local
###    nameserver[0] : 10.5.0.18
###    nameserver[1] : 10.5.0.19
###    nameserver[2] : 192.168.7.85
###    if_index : 17 (utun2)
###    flags    : Scoped, Request A records, Request AAAA records
###    reach    : 0x00000002 (Reachable)
###    order    : 1


###  $ printf 'list ".*DNS"' | scutil
###    subKey [0] = Setup:/Network/Service/8FA97706-CAE9-4003-93A4-E7F7216CCF0B/DNS
###    subKey [1] = State:/Network/Global/DNS
###    subKey [2] = State:/Network/MulticastDNS
###    subKey [3] = State:/Network/PrivateDNS
###    subKey [4] = State:/Network/Service/8FA97706-CAE9-4003-93A4-E7F7216CCF0B/DNS
###    subKey [5] = State:/Network/Service/com.cisco.anyconnect/DNS

### https://bandwidth.slack.com/archives/C01A01GQ1JM/p1620317818186000



###  $ printf "get State:/Network/Service/com.cisco.anyconnect/IPv4\nd.show\n" | scutil | grep Interface
###    InterfaceName : utun2
