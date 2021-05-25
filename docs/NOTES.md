```
func scutilCmd() string {
	cmd := `printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`
	out, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	output := string(out[:])
	return output
}
```

# off vpn
```
$ printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil
  No such key
<dictionary> {
}
```

# on vpn
```
$ printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil
<dictionary> {
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
}
```
