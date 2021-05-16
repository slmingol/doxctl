# doxctl

## TLDR
`doxctl` is a diagnostic CLI tool that endusers can use to triage connectivity problems. It can help with the following areas:

- DNS - resolvers, search paths, etc. are set
- Resolvers - DNS resolvers defined and reachable

## Example Output
```
$ doxctl dns diag
dns diag
model_cmds/02_dns.sh

DNS Resolver Checks
===================


DomainName       set
SearchDomains    set
ServerAddresses  set



Ping Resolver Checks
====================


How many resolvers found? 		 ---> 2 <---

Was resolver 10.5.0.18 pingable? 	 ---> yes <---
Can we reach port 53 (DNS) via TCP? 	 ---> yes <--- 		 [ Connection to 10.5.0.18 port 53 [tcp/domain] succeeded! ]
Can we reach port 53 (DNS) via UDP? 	 ---> yes <--- 		 [ Connection to 10.5.0.18 port 53 [udp/domain] succeeded! ]


Was resolver 10.5.0.19 pingable? 	 ---> yes <---
Can we reach port 53 (DNS) via TCP? 	 ---> yes <--- 		 [ Connection to 10.5.0.19 port 53 [tcp/domain] succeeded! ]
Can we reach port 53 (DNS) via UDP? 	 ---> yes <--- 		 [ Connection to 10.5.0.19 port 53 [udp/domain] succeeded! ]



Dig Resolver Checks
===================


	resolver: 10.5.0.18 | site: lab1 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.19 | site: lab1 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.18 | site: rdu1 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.19 | site: rdu1 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.18 | site: atl1 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.19 | site: atl1 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.18 | site: dfw1 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.19 | site: dfw1 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.18 | site: lax2 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.19 | site: lax2 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.18 | site: jfk1 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |
	resolver: 10.5.0.19 | site: jfk1 | idm-01a | cnt: 1 | idm-01b |  cnt: 1 |


Check if we can resolve all 24 IDM server names (idm-01[ab].*)? 	 ---> yes <--- 		 [ Actual: 24 ]


```
