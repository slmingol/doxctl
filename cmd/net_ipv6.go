package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"doxctl/internal/output"
)

type ipv6HostResult struct {
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
	Host      string    `json:"host" yaml:"host"`
	IPv6Addrs []string  `json:"ipv6Addrs" yaml:"ipv6Addrs"`
	HasAAAA   bool      `json:"hasAAAA" yaml:"hasAAAA"`
	HasA      bool      `json:"hasA" yaml:"hasA"`
	DualStack bool      `json:"dualStack" yaml:"dualStack"`
	Reachable bool      `json:"reachable" yaml:"reachable"`
	LatencyMs float64   `json:"latencyMs" yaml:"latencyMs"`
	Error     string    `json:"error,omitempty" yaml:"error,omitempty"`
}

type ipv6RouteInfo struct {
	HasDefaultRoute bool     `json:"hasDefaultRoute" yaml:"hasDefaultRoute"`
	RouteCount      int      `json:"routeCount" yaml:"routeCount"`
	Routes          []string `json:"routes" yaml:"routes"`
}

type netIPv6Output struct {
	Timestamp   time.Time        `json:"timestamp" yaml:"timestamp"`
	HostResults []ipv6HostResult `json:"hostResults" yaml:"hostResults"`
	RouteInfo   ipv6RouteInfo    `json:"routeInfo" yaml:"routeInfo"`
	Summary     struct {
		TotalHosts  int `json:"totalHosts" yaml:"totalHosts"`
		DualStack   int `json:"dualStack" yaml:"dualStack"`
		IPv6Only    int `json:"ipv6Only" yaml:"ipv6Only"`
		Reachable   int `json:"reachable" yaml:"reachable"`
		Unreachable int `json:"unreachable" yaml:"unreachable"`
	} `json:"summary" yaml:"summary"`
}

func netIPv6Check() {
	netIPv6CheckWithDeps(conf, netPacketCount, NewPinger, NewCommandExecutor())
}

func netIPv6CheckWithDeps(config *config, packetCount int, pingerFactory func(string) (Pinger, error), executor CommandExecutor) {
	result := netIPv6Output{
		Timestamp:   time.Now(),
		HostResults: []ipv6HostResult{},
	}

	var targets []string
	expander := NewBraceExpander()
	hostMap := make(map[string]bool)
	for _, service := range config.Svcs {
		for _, server := range service.Svrs {
			for _, h := range expander.Expand(server) {
				if !hostMap[h] {
					hostMap[h] = true
					targets = append(targets, h)
				}
			}
		}
	}

	if len(targets) == 0 {
		fmt.Println("")
		fmt.Printf("\033[1;33mWARNING:\033[0m No targets configured in services\n")
		return
	}

	result.RouteInfo = getIPv6Routes(executor)

	_ = RunWithSpinnerProgress("Checking IPv6 connectivity", len(targets), func(index int) error {
		host := targets[index]
		hr := ipv6HostResult{
			Timestamp: time.Now(),
			Host:      host,
			IPv6Addrs: []string{},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
		if err != nil {
			hr.Error = err.Error()
			result.HostResults = append(result.HostResults, hr)
			result.Summary.Unreachable++
			return nil
		}

		for _, addr := range addrs {
			ip := addr.IP
			if ip.To4() == nil {
				hr.HasAAAA = true
				hr.IPv6Addrs = append(hr.IPv6Addrs, ip.String())
			} else {
				hr.HasA = true
			}
		}

		hr.DualStack = hr.HasA && hr.HasAAAA

		if hr.HasAAAA {
			pingTarget := hr.IPv6Addrs[0]
			pinger, perr := pingerFactory(pingTarget)
			if perr == nil {
				pinger.SetCount(packetCount)
				pinger.SetTimeout(10 * time.Second)
				if rerr := pinger.Run(); rerr == nil {
					stats := pinger.Statistics()
					if stats.PacketsRecv > 0 {
						hr.Reachable = true
						hr.LatencyMs = float64(stats.AvgRtt.Microseconds()) / 1000.0
					}
				}
			}
		}

		if hr.Reachable {
			result.Summary.Reachable++
		} else {
			result.Summary.Unreachable++
		}
		if hr.DualStack {
			result.Summary.DualStack++
		} else if hr.HasAAAA && !hr.HasA {
			result.Summary.IPv6Only++
		}

		result.HostResults = append(result.HostResults, hr)
		return nil
	})

	result.Summary.TotalHosts = len(result.HostResults)

	switch outputFormat {
	case "json":
		output.PrintJSON(result)
	case "yaml":
		output.PrintYAML(result)
	default:
		printNetIPv6Table(result)
	}
}

func getIPv6Routes(executor CommandExecutor) ipv6RouteInfo {
	info := ipv6RouteInfo{Routes: []string{}}

	var out []byte
	var err error
	if runtime.GOOS == "darwin" {
		out, err = executor.Execute("netstat", "-rn", "-f", "inet6")
	} else {
		out, err = executor.Execute("ip", "-6", "route", "show")
	}
	if err != nil {
		return info
	}

	for line := range strings.SplitSeq(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "default") || strings.HasPrefix(line, "::/0") {
			info.HasDefaultRoute = true
		}
		if strings.Contains(line, ":") || strings.HasPrefix(line, "default") {
			if info.RouteCount < 20 {
				info.Routes = append(info.Routes, line)
			}
			info.RouteCount++
		}
	}

	return info
}

func printNetIPv6Table(result netIPv6Output) {
	headers := []string{"Host", "IPv6 Addresses", "AAAA", "Dual-Stack", "Reachable", "Latency (ms)"}
	var rows [][]string

	boolMark := func(b bool) string {
		if b {
			return "✓"
		}
		return "✗"
	}

	for _, hr := range result.HostResults {
		addrs := strings.Join(hr.IPv6Addrs, ", ")
		if addrs == "" {
			addrs = "-"
		}
		latency := "-"
		if hr.Reachable {
			latency = fmt.Sprintf("%.2f", hr.LatencyMs)
		}
		rows = append(rows, []string{
			hr.Host,
			addrs,
			boolMark(hr.HasAAAA),
			boolMark(hr.DualStack),
			boolMark(hr.Reachable),
			latency,
		})
	}

	fmt.Print(createStyledTableWithTypedSeparators(headers, rows, "IPv6 Connectivity Check", nil))

	fmt.Printf("\nSummary: %d hosts | %d dual-stack | %d IPv6-only | %d reachable\n",
		result.Summary.TotalHosts,
		result.Summary.DualStack,
		result.Summary.IPv6Only,
		result.Summary.Reachable,
	)

	defaultRoute := "✗"
	if result.RouteInfo.HasDefaultRoute {
		defaultRoute = "✓"
	}
	fmt.Printf("IPv6 routing: default route %s | %d routes found\n\n", defaultRoute, result.RouteInfo.RouteCount)

	if len(result.HostResults) > 0 {
		var errs []ipv6HostResult
		for _, hr := range result.HostResults {
			if hr.Error != "" {
				errs = append(errs, hr)
			}
		}
		if len(errs) > 0 {
			fmt.Printf("\033[1;33mLookup errors:\033[0m\n")
			for _, hr := range errs {
				fmt.Printf("  • %s: %s\n", hr.Host, hr.Error)
			}
			fmt.Println()
		}
	}

	if len(result.HostResults) == 0 {
		fmt.Printf("\033[1;31mERROR:\033[0m No results collected\n\n")
		os.Exit(1)
	}
}
