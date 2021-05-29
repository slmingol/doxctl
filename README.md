![dox logo](https://github.com/slmingol/doxctl/blob/main/imgs/whats_up_dox__banner.png?raw=true)

## TLDR
`doxctl` is a diagnostic CLI tool that endusers can use to triage connectivity problems stemming from their VPN & DNS setups on their laptops. It can help with the following areas:

| Area | Description |
| ---- | ----------- |
| VPN  | Can servers be reached over the VPN across geo-locations |
| DNS  | Resolvers, search paths, etc. are set |
| Resolvers | VPN defined DNS resolvers are defined and reachable |
| Routing and connectivity | When on the VPN well-known servers are reachable |

## Usage
### General
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

### DNS
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

### DNS Example Output

#### resolverChk

##### Off VPN
<details><summary>Tree - CLICK ME</summary>
<p>

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
</p>
</details>

##### On VPN
<details><summary>Tree - CLICK ME</summary>
<p>

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

##### Off VPN
<details><summary>Tree - CLICK ME</summary>
<p>
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
</p>
</details>

##### On VPN
<details><summary>Tree - CLICK ME</summary>
<p>
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

