# doxctl

## TLDR
`doxctl` is a diagnostic CLI tool that endusers can use to triage connectivity problems stemming from their VPN & DNS setups on their laptops. It can help with the following areas:

| Area | Description |
| ---- | ----------- |
| VPN  | Can servers be reached over the VPN across geo-locations |
| DNS  | Resolvers, search paths, etc. are set |
| Resolvers | VPN defined DNS resolvers are defined and reachable |
| Routing and connectivity | When on the VPN well-known servers are reachable |

<img src="./imgs/0bbb4991cac283330a29c537711f0ac2_whats_up_doc.jpg" style="width:1000px;height:200px;">

## Usage
### General
```
$ doxctl -h

'doxctl' is a collection of tools which can be used to diagnose:

  - DNS resolvers
  - VPN access
  - General access to well-known servers
  - ... or general network connectivity issues

Usage:
  doxctl [flags]
  doxctl [command]

Available Commands:
  dns         Run diagnostics related to DNS servers' (resolvers') configurations
  help        Help about any command
  net         A brief description of your command
  svcs        A brief description of your command
  svrs        A brief description of your command
  vpn         A brief description of your command

Flags:
  -h, --help      help for doxctl
  -t, --toggle    Help message for toggle
  -v, --verbose   Enable verbose output of commands

Use "doxctl [command] --help" for more information about a command.


```

### DNS
```
$ doxctl dns -h

doxctl's 'dns' subcommand can help triage DNS resovler configuration issues,
general access to DNS resolvers and name resolution against DNS resolvers.

Usage:
  doxctl dns [flags]

Flags:
  -d, --digChk        check if VPN defined resolvers respond with well-known servers in DCs
  -h, --help          help for dns
  -p, --pingChk       check if VPN defined resolvers are pingable & reachable
  -r, --resolverChk   check if VPN designated DNS resolvers are configured

Global Flags:
  -v, --verbose   Enable verbose output of commands


```

## DNS Example Output

#### resolverChk

<details><summary>Tree - CLICK ME</summary>
<p>

```
$ doxctl dns -r


DNS Resolver Checks
===================


DomainName       set
SearchDomains    set
ServerAddresses  set




```
</p>
</details>

#### pingChk

<details><summary>Tree - CLICK ME</summary>
<p>

```
$ doxctl dns -p


Ping Resolver Checks
====================


How many resolvers found? 		 ---> 2 <---

Was resolver 10.5.0.18 pingable? 	 ---> yes <---
Can we reach port 53 (DNS) via TCP? 	 ---> yes <--- 		 [ Connection to 10.5.0.18 port 53 [tcp/domain] succeeded! ]
Can we reach port 53 (DNS) via UDP? 	 ---> yes <--- 		 [ Connection to 10.5.0.18 port 53 [udp/domain] succeeded! ]


Was resolver 10.5.0.19 pingable? 	 ---> yes <---
Can we reach port 53 (DNS) via TCP? 	 ---> yes <--- 		 [ Connection to 10.5.0.19 port 53 [tcp/domain] succeeded! ]
Can we reach port 53 (DNS) via UDP? 	 ---> yes <--- 		 [ Connection to 10.5.0.19 port 53 [udp/domain] succeeded! ]





```
</p>
</details>

#### digChk

<details><summary>Tree - CLICK ME</summary>
<p>

```
$ doxctl dns -d


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
</p>
</details>

#### Debugging
All the CLI subcommands can make use of either the `-v` or the `--verbose` switch to gather further diagnostics which can be helpful when triaging connectivity issues.

<details><summary>Tree - CLICK ME</summary>
<p>

For example:
```
$ doxctl dns -r -v
+ printf '\n\nDNS Resolver Checks\n===================\n\n\n'


DNS Resolver Checks
===================


++ printf 'get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n'
++ scutil
+ vpn_resolvers='<dictionary> {
  DomainName : bandwidth.local
  SearchDomains : <array> {
    0 : bandwidth.local
  }
  SearchOrder : 1
  ServerAddresses : <array> {
    0 : 10.5.0.18
    1 : 10.5.0.19
    2 : 192.168.7.85
  }
  SupplementalMatchDomains : <array> {
    0 :
    1 : bandwidth.local
  }
}'
+ column -t
+ echo '<dictionary> {
  DomainName : bandwidth.local
  SearchDomains : <array> {
    0 : bandwidth.local
  }
  SearchOrder : 1
  ServerAddresses : <array> {
    0 : 10.5.0.18
    1 : 10.5.0.19
    2 : 192.168.7.85
  }
  SupplementalMatchDomains : <array> {
    0 :
    1 : bandwidth.local
  }
}'
+ grep -q 'DomainName.*bandwidth.local'
+ echo 'DomainName set'
+ echo '<dictionary> {
  DomainName : bandwidth.local
  SearchDomains : <array> {
    0 : bandwidth.local
  }
  SearchOrder : 1
  ServerAddresses : <array> {
    0 : 10.5.0.18
    1 : 10.5.0.19
    2 : 192.168.7.85
  }
  SupplementalMatchDomains : <array> {
    0 :
    1 : bandwidth.local
  }
}'
+ grep -A1 SearchDomains
+ grep -qE '[0-1].*bandwidth'
+ echo 'SearchDomains set'
+ echo '<dictionary> {
  DomainName : bandwidth.local
  SearchDomains : <array> {
    0 : bandwidth.local
  }
  SearchOrder : 1
  ServerAddresses : <array> {
    0 : 10.5.0.18
    1 : 10.5.0.19
    2 : 192.168.7.85
  }
  SupplementalMatchDomains : <array> {
    0 :
    1 : bandwidth.local
  }
}'
+ grep -A3 ServerAddresses
+ grep -qE '[0-1].*10.5'
+ echo 'ServerAddresses set'
DomainName       set
SearchDomains    set
ServerAddresses  set
+ [[ 1 -eq 1 ]]
+ set +x




```
</p>
</details>

### RESOURCES

#### References
- [Exploring Go Packages: Cobra](https://levelup.gitconnected.com/exploring-go-packages-cobra-fce6c4e331d6)
- [CLI Command Line SDK - Cobra](https://github.com/spf13/cobra)
- [Building a multipurpose CLI tool with Cobra and Go](https://dev.to/lumexralph/building-a-multipurpose-cli-tool-with-cobra-and-go-2492)
- [How to create a CLI in golang with cobra](https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177)

#### Example CLI tools written in Go
- [docker/hub-tool](ttps://github.com/docker/hub-tool/tree/main/internal/commands)
