![dox logo](https://github.com/slmingol/doxctl/blob/main/imgs/whats_up_dox__banner.png?raw=true)

# TLDR
`doxctl` is a diagnostic CLI tool that endusers can use to triage connectivity problems stemming from their VPN & DNS setups on their laptops. It can help with the following areas:

| Area | Description |
| ---- | ----------- |
| VPN  | Can servers be reached over the VPN across geo-locations |
| DNS  | Resolvers, search paths, etc. are set |
| Resolvers | VPN defined DNS resolvers are defined and reachable |
| Routing and connectivity | When on the VPN well-known servers are reachable |

# Installation
## MacOS
<details><summary>Tree - CLICK ME</summary>
<p>

```
$ brew install slmingol/tap/doxctl
```

### Examples
#### Install
```
$ brew install slmingol/tap/doxctl
==> Installing doxctl from slmingol/tap
==> Downloading https://ghcr.io/v2/homebrew/core/go/manifests/1.16.4
######################################################################## 100.0%
==> Downloading https://ghcr.io/v2/homebrew/core/go/blobs/sha256:8aa23b394e05aaef495604670401f9308c01c3c30a1857077493a80d1719e089
==> Downloading from https://pkg-containers.githubusercontent.com/ghcr1/blobs/sha256:8aa23b394e05aaef495604670401f9308c01c3c30a1857077493a80d1719e089?se=202
######################################################################## 100.0%
==> Downloading https://github.com/slmingol/doxctl/releases/download/0.0.27-alpha/doxctl_0.0.27-alpha_Darwin_x86_64.tar.gz
==> Downloading from https://github-releases.githubusercontent.com/367779289/09fcfe00-c133-11eb-9f76-3f053be7f72c?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Cre
######################################################################## 100.0%
==> Installing dependencies for slmingol/tap/doxctl: go
==> Installing slmingol/tap/doxctl dependency: go
==> Pouring go--1.16.4.catalina.bottle.tar.gz
🍺  /usr/local/Cellar/go/1.16.4: 9,956 files, 503.6MB
==> Installing slmingol/tap/doxctl
🍺  /usr/local/Cellar/doxctl/0.0.27-alpha: 5 files, 9.1MB, built in 4 seconds
```

#### search
```
$ brew search doxctl
==> Formulae
slmingol/tap/doxctl ✔
```

#### Upgrade
```
$ brew update
Updated 2 taps (homebrew/cask and slmingol/tap).
==> Updated Formulae
slmingol/tap/doxctl ✔
==> Deleted Casks
appstudio                              filedrop                               hex                                    rss

You have 96 outdated formulae and 10 outdated casks installed.
You can upgrade them with brew upgrade
or list them with brew outdated.


$ brew upgrade doxctl
==> Upgrading 1 outdated package:
slmingol/tap/doxctl 0.0.27-alpha -> 0.0.28-alpha
==> Upgrading slmingol/tap/doxctl 0.0.27-alpha -> 0.0.28-alpha
==> Downloading https://github.com/slmingol/doxctl/releases/download/0.0.28-alpha/doxctl_0.0.28-alpha_Darwin_x86_64.tar.gz
==> Downloading from https://github-releases.githubusercontent.com/367779289/c7d4bc00-c134-11eb-8b3d-622ae425b081?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Cre
######################################################################## 100.0%
🍺  /usr/local/Cellar/doxctl/0.0.28-alpha: 5 files, 9.1MB, built in 4 seconds
Removing: /usr/local/Cellar/doxctl/0.0.27-alpha... (5 files, 9.1MB)
Removing: /Users/smingolelli/Library/Caches/Homebrew/doxctl--0.0.27-alpha.tar.gz... (3.4MB)
```

#### Uninstall
```
$ brew uninstall doxctl
Uninstalling /usr/local/Cellar/doxctl/0.0.27-alpha... (5 files, 9.1MB)
```
</p>
</details>

# Usage
## General
```
$ doxctl -h

'doxctl' is a collection of tools which can be used to diagnose & triage problems
stemming from the following areas with a laptop or desktop system:

  - DNS, specifically with the configuration of resolvers
  - VPN configuration and network connectivity over it
  - General access to well-known servers
  - General access to well-known services
  - ... or general network connectivity issues

Usage:
  doxctl [command]

Available Commands:
  dns         Run diagnostics related to DNS servers' (resolvers') configurations
  help        Help about any command
  net         TBD
  svcs        TBD
  svrs        TBD
  vpn         Run diagnostics related to VPN connections, net i/fs & configurations

Flags:
  -c, --config string   config file (default is $HOME/.doxctl.yaml)
  -h, --help            help for doxctl
  -v, --verbose         Enable verbose output of commands

Use "doxctl [command] --help" for more information about a command.
```

------------------------------------------------------------------------------

## DNS
```
$ doxctl dns -h

doxctl's 'dns' subcommand can help triage DNS resovler configuration issues,
general access to DNS resolvers and name resolution against DNS resolvers.

Usage:
  doxctl dns [flags]

Flags:
  -a, --allChk        Run all the checks in this subcommand module
  -d, --digChk        Check if VPN defined resolvers respond with well-known servers in DCs
  -h, --help          help for dns
  -p, --pingChk       Check if VPN defined resolvers are pingable & reachable
  -r, --resolverChk   Check if VPN designated DNS resolvers are configured

Global Flags:
  -c, --config string   config file (default is $HOME/.doxctl.yaml)
  -v, --verbose         Enable verbose output of commands
```

### Example Output

#### resolverChk

<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl dns -r

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml


┌───────────────────────────────────────────────────────────────────────────┐
│ VPN defined DNS Resolver Checks                                           │
├──────────────────────────────────────────┬────────────────────────────────┤
│ PROPERTY DESCRIPTION                     │ VALUE                          │
├──────────────────────────────────────────┼────────────────────────────────┤
│ DomainName defined?                      │ unset                          │
│ SearchDomains defined?                   │ unset                          │
│ ServerAddresses defined?                 │ unset                          │
└──────────────────────────────────────────┴────────────────────────────────┘

** NOTE:** Any values of unset indicate that the VPN client is not defining DNS resolver(s) properly!


```

##### On VPN
```
$ doxctl dns -r

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml


┌───────────────────────────────────────────────────────────────────────────┐
│ VPN defined DNS Resolver Checks                                           │
├──────────────────────────────────────────┬────────────────────────────────┤
│ PROPERTY DESCRIPTION                     │ VALUE                          │
├──────────────────────────────────────────┼────────────────────────────────┤
│ DomainName defined?                      │ set                            │
│ SearchDomains defined?                   │ set                            │
│ ServerAddresses defined?                 │ set                            │
└──────────────────────────────────────────┴────────────────────────────────┘

** NOTE:** Any values of unset indicate that the VPN client is not defining DNS resolver(s) properly!


```
</p>
</details>

#### pingChk

<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl dns -p

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml

┌──────────────────────────────────────────────────────────────────────────────────┐
│ VPN defined DNS Resolver Connectivity Checks                                     │
├──────────────────────────────────────────┬───────────────┬───────────────┬───────┤
│                     PROPERTY DESCRIPTION │            IP │       NET I/F │ VALUE │
├──────────────────────────────────────────┼───────────────┼───────────────┼───────┤
└──────────────────────────────────────────┴───────────────┴───────────────┴───────┘

** WARN:** Your VPN client does not appear to be defining any DNS resolver(s) properly,
           you're either not connected via VPN or it's misconfigured!



```

##### On VPN
```
$ doxctl dns -p

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml

┌──────────────────────────────────────────────────────────────────────────────────┐
│ VPN defined DNS Resolver Connectivity Checks                                     │
├──────────────────────────────────────────┬───────────────┬───────────────┬───────┤
│ PROPERTY DESCRIPTION                     │ IP            │ NET I/F       │ VALUE │
├──────────────────────────────────────────┼───────────────┼───────────────┼───────┤
│ Resovler is pingable?                    │ 10.5.0.18     │ utun2         │ true  │
│ Reachable via TCP?                       │ 10.5.0.18     │ utun2         │ true  │
│ Reachable via UDP?                       │ 10.5.0.18     │ utun2         │ true  │
├──────────────────────────────────────────┼───────────────┼───────────────┼───────┤
│ Resovler is pingable?                    │ 10.5.0.19     │ utun2         │ true  │
│ Reachable via TCP?                       │ 10.5.0.19     │ utun2         │ true  │
│ Reachable via UDP?                       │ 10.5.0.19     │ utun2         │ true  │
└──────────────────────────────────────────┴───────────────┴───────────────┴───────┘



```
</p>
</details>

#### digChk

<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl dns -d

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml

┌──────────────────────────────────────────────────────────────────────────────┐
│ Dig Check against VPN defined DNS Resolvers                                  │
├──────────────────────────────────────────┬─────────────────┬─────────────────┤
│ HOSTNAME TO 'DIG'                        │ RESOLVER IP     │ IS RESOLVABLE?  │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.lab1.somedom.local               │                 │ false           │
│ idm-01b.lab1.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.rdu1.somedom.local               │                 │ false           │
│ idm-01b.rdu1.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.atl1.somedom.local               │                 │ false           │
│ idm-01b.atl1.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.dfw1.somedom.local               │                 │ false           │
│ idm-01b.dfw1.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.lax2.somedom.local               │                 │ false           │
│ idm-01b.lax2.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.jfk1.somedom.local               │                 │ false           │
│ idm-01b.jfk1.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ SUCCESSESFUL QUERIES                     │ RESOLVER #1: 0  │                 │
│                                          │ RESOLVER #2: 0  │                 │
└──────────────────────────────────────────┴─────────────────┴─────────────────┘

** WARN:** Your VPN client does not appear to be defining any DNS resolver(s) properly,
           you're either not connected via VPN or it's misconfigured!



```

##### On VPN
```
$ doxctl dns -d

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml

┌──────────────────────────────────────────────────────────────────────────────┐
│ Dig Check against VPN defined DNS Resolvers                                  │
├──────────────────────────────────────────┬─────────────────┬─────────────────┤
│ HOSTNAME TO 'DIG'                        │ RESOLVER IP     │ IS RESOLVABLE?  │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.lab1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.lab1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.lab1.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.lab1.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.rdu1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.rdu1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.rdu1.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.rdu1.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.atl1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.atl1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.atl1.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.atl1.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.dfw1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.dfw1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.dfw1.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.dfw1.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.lax2.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.lax2.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.lax2.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.lax2.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.jfk1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.jfk1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.jfk1.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.jfk1.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ SUCCESSESFUL QUERIES                     │ RESOLVER #1: 12 │                 │
│                                          │ RESOLVER #2: 12 │                 │
└──────────────────────────────────────────┴─────────────────┴─────────────────┘



```
</p>
</details>

------------------------------------------------------------------------------

## VPN
```
$ doxctl vpn -h

doxctl's 'vpn' subcommand can help triage VPN related configuration issues,
& routes related to a VPN connection.

Usage:
  doxctl vpn [flags]

Flags:
  -a, --allChk           Run all the checks in this subcommand module
  -h, --help             help for vpn
  -i, --ifReachableChk   Check if network interfaces are reachable
  -r, --vpnRoutesChk     Check if >5 VPN routes are defined
  -s, --vpnStatusChk     Check if VPN client's status reports as 'Connected'

Global Flags:
  -c, --config string   config file (default is $HOME/.doxctl.yaml)
  -v, --verbose         Enable verbose output of commands
```

### Example Output

#### ifReachableChk
<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl vpn -i

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│ Interfaces Reachable Checks                                                                 │
├────────────────────────────────────────────────────┬────────────────────────────────┬───────┤
│ PROPERTY DESCRIPTION                               │ VALUE                          │ NOTES │
├────────────────────────────────────────────────────┼────────────────────────────────┼───────┤
│ How many network interfaces found?                 │ 1                              │ [en0] │
│ At least 1 interface's a utun device?              │ false                          │ []    │
│ All active interfaces are reporting as reachable?  │ true                           │       │
└────────────────────────────────────────────────────┴────────────────────────────────┴───────┘

WARNING: Your VPN client does not appear to be defining a TUN interface properly,
WARNING: you're VPN is either not connected or it's misconfigured!



```

##### On VPN
```
$ doxctl vpn -i

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌───────────────────────────────────────────────────────────────────────────────────────────────────┐
│ Interfaces Reachable Checks                                                                       │
├────────────────────────────────────────────────────┬────────────────────────────────┬─────────────┤
│ PROPERTY DESCRIPTION                               │ VALUE                          │ NOTES       │
├────────────────────────────────────────────────────┼────────────────────────────────┼─────────────┤
│ How many network interfaces found?                 │ 2                              │ [en0 utun2] │
│ At least 1 interface's a utun device?              │ true                           │ [utun2]     │
│ All active interfaces are reporting as reachable?  │ true                           │             │
└────────────────────────────────────────────────────┴────────────────────────────────┴─────────────┘



```
</p>
</details>

#### vpnRoutesChk
<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl vpn -r

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│ VPN Interface Route Checks                                                                  │
├────────────────────────────────────────────────────┬────────────────────────────────┬───────┤
│ PROPERTY DESCRIPTION                               │ VALUE                          │ NOTES │
├────────────────────────────────────────────────────┼────────────────────────────────┼───────┤
│ At least 5 routes using interface [NIL]?           │ false                          │     0 │
└────────────────────────────────────────────────────┴────────────────────────────────┴───────┘

WARNING: Your VPN client does not appear to be defining a TUN interface properly,
WARNING: you're VPN is either not connected or it's misconfigured!



```

##### On VPN
```
$ doxctl vpn -r

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│ VPN Interface Route Checks                                                                  │
├────────────────────────────────────────────────────┬────────────────────────────────┬───────┤
│ PROPERTY DESCRIPTION                               │ VALUE                          │ NOTES │
├────────────────────────────────────────────────────┼────────────────────────────────┼───────┤
│ At least 5 routes using interface [utun2]?         │ true                           │   148 │
└────────────────────────────────────────────────────┴────────────────────────────────┴───────┘



```
</p>
</details>

#### vpnStatusChk
<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl vpn -s

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌───────────────────────────────────────────────────────────────────────────────────────────────┐
│ VPN Connection Status Checks                                                                  │
├──────────────────────────────────────────────────────┬────────────────────────────────┬───────┤
│ PROPERTY DESCRIPTION                                 │ VALUE                          │ NOTES │
├──────────────────────────────────────────────────────┼────────────────────────────────┼───────┤
│ VPN Client reports connection status as 'Connected'? │ false                          │       │
└──────────────────────────────────────────────────────┴────────────────────────────────┴───────┘

WARNING: Your VPN client's does not appear to be a state of 'connected',
WARNING: it's either down or misconfigured!



```

##### On VPN
```
$ doxctl vpn -s

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌───────────────────────────────────────────────────────────────────────────────────────────────┐
│ VPN Connection Status Checks                                                                  │
├──────────────────────────────────────────────────────┬────────────────────────────────┬───────┤
│ PROPERTY DESCRIPTION                                 │ VALUE                          │ NOTES │
├──────────────────────────────────────────────────────┼────────────────────────────────┼───────┤
│ VPN Client reports connection status as 'Connected'? │ true                           │       │
└──────────────────────────────────────────────────────┴────────────────────────────────┴───────┘



```
</p>
</details>

------------------------------------------------------------------------------

## Debugging - TBD/WIP
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

## Resources

### References
- [Exploring Go Packages: Cobra](https://levelup.gitconnected.com/exploring-go-packages-cobra-fce6c4e331d6)
- [CLI Command Line SDK - Cobra](https://github.com/spf13/cobra)
- [Building a multipurpose CLI tool with Cobra and Go](https://dev.to/lumexralph/building-a-multipurpose-cli-tool-with-cobra-and-go-2492)
- [How to create a CLI in golang with cobra](https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177)
- [How to pipe several commands in Go?](https://stackoverflow.com/questions/10781516/how-to-pipe-several-commands-in-go)
- [Executing System Commands With Golang](https://tutorialedge.net/golang/executing-system-commands-with-golang/)
- [gookit/color](https://github.com/gookit/color)

### Example CLI tools written in Go
- [docker/hub-tool](https://github.com/docker/hub-tool/tree/main/internal/commands)

### MacOS
- [SCNetworkReachability](https://developer.apple.com/documentation/systemconfiguration/scnetworkreachability-g7d)
- [scutil generalized interface to "dynamic Store" and Network Services](https://www.real-world-systems.com/docs/scutil.1.html)
- [Homebrew formulas - slmingol/homebrew-tap](https://github.com/slmingol/homebrew-tap)
