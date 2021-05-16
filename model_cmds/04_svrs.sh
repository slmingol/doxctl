#!/usr/bin/env bash

# $ dscacheutil -q host -a name ocp-app-01a.lab1.examplecorp.local
# name: ocp-app-01a.lab1.examplecorp.local
# ip_address: 192.168.113.138


### MacOS


# check for DNS lookup results for well known servers

# lab1, rdu1, atl1, dfw1, jfk1, lax2
for site in lab1 rdu1 dfw1 jfk1 lax2; do
    dscacheutil -q host -a name ocp-master-01a.${site}.examplecorp.local
    dscacheutil -q host -a name ocp-app-01a.${site}.examplecorp.local
    dscacheutil -q host -a name idm-01a.${site}.examplecorp.local
    dscacheutil -q host -a name idm-01b.${site}.examplecorp.local
done

# $ printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil
# <dictionary> {
#   DomainName : example.local
#   SearchDomains : <array> {
#     0 : example.local
#   }
#   SearchOrder : 1
#   ServerAddresses : <array> {
#     0 : 10.5.0.18
#     1 : 10.5.0.19
#     2 : 192.168.7.85
#   }
#   SupplementalMatchDomains : <array> {
#     0 :
#     1 : example.local
#   }
# }


