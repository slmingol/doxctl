package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/go-ping/ping"
)

func makeIPv6Config(servers []string) *config {
	return &config{
		DomainName: "example.com",
		Sites:      []string{"site1"},
		Svcs:       []svc{{Svc: "test", Svrs: servers}},
	}
}

func ipv6PingerOK() func(string) (Pinger, error) {
	return func(addr string) (Pinger, error) {
		return &mockNetPinger{
			stats: &ping.Statistics{
				PacketsRecv: 10,
				PacketsSent: 10,
				PacketLoss:  0.0,
				AvgRtt:      15 * time.Millisecond,
			},
		}, nil
	}
}

func ipv6PingerFail() func(string) (Pinger, error) {
	return func(addr string) (Pinger, error) {
		return nil, errors.New("no route to host")
	}
}

func TestNetIPv6CheckWithDeps_NoTargets(t *testing.T) {
	cfg := &config{DomainName: "example.com", Sites: []string{}, Svcs: []svc{}}
	exec := &mockCommandExecutor{commands: map[string][]byte{}}
	// Should return without panic
	netIPv6CheckWithDeps(cfg, 5, ipv6PingerOK(), exec)
}

func TestNetIPv6CheckWithDeps_PingerFactoryError(t *testing.T) {
	cfg := makeIPv6Config([]string{"host1.example.com"})
	exec := &mockCommandExecutor{commands: map[string][]byte{}}
	// Pinger fails -- should not panic
	netIPv6CheckWithDeps(cfg, 5, ipv6PingerFail(), exec)
}

func TestNetIPv6CheckWithDeps_MultipleHosts(t *testing.T) {
	cfg := makeIPv6Config([]string{"host1.example.com", "host2.example.com", "host3.example.com"})
	exec := &mockCommandExecutor{commands: map[string][]byte{}}
	netIPv6CheckWithDeps(cfg, 5, ipv6PingerOK(), exec)
}

func TestNetIPv6CheckWithDeps_BraceExpansion(t *testing.T) {
	cfg := makeIPv6Config([]string{"host{1,2,3}.example.com"})
	exec := &mockCommandExecutor{commands: map[string][]byte{}}
	netIPv6CheckWithDeps(cfg, 5, ipv6PingerOK(), exec)
}

func TestGetIPv6Routes_DefaultRoute(t *testing.T) {
	routeOutput := `default via fe80::1 dev eth0
2001:db8::/32 dev eth0 proto kernel metric 256
fe80::/64 dev eth0 proto kernel metric 256
`
	// Use execFunc to return output regardless of platform command
	exec := &mockCommandExecutor{
		execFunc: func(name string, args ...string) ([]byte, error) {
			return []byte(routeOutput), nil
		},
	}
	info := getIPv6Routes(exec)
	if !info.HasDefaultRoute {
		t.Error("expected HasDefaultRoute=true")
	}
	if info.RouteCount == 0 {
		t.Error("expected RouteCount > 0")
	}
}

func TestGetIPv6Routes_CIDRDefaultRoute(t *testing.T) {
	routeOutput := `::/0 via fe80::1 dev eth0 proto ra metric 100
2001:db8::/32 dev eth0 proto kernel metric 256
`
	exec := &mockCommandExecutor{
		execFunc: func(name string, args ...string) ([]byte, error) {
			return []byte(routeOutput), nil
		},
	}
	info := getIPv6Routes(exec)
	if !info.HasDefaultRoute {
		t.Error("expected HasDefaultRoute=true for ::/0")
	}
	if info.RouteCount == 0 {
		t.Error("expected RouteCount > 0")
	}
}

func TestGetIPv6Routes_ExecutorError(t *testing.T) {
	exec := &mockCommandExecutor{err: errors.New("command not found")}
	info := getIPv6Routes(exec)
	if info.HasDefaultRoute {
		t.Error("expected HasDefaultRoute=false on error")
	}
	if info.RouteCount != 0 {
		t.Error("expected RouteCount=0 on error")
	}
}

func TestGetIPv6Routes_NoDefaultRoute(t *testing.T) {
	exec := &mockCommandExecutor{
		execFunc: func(name string, args ...string) ([]byte, error) {
			return []byte("fe80::/64 dev eth0 proto kernel metric 256\n"), nil
		},
	}
	info := getIPv6Routes(exec)
	if info.HasDefaultRoute {
		t.Error("expected HasDefaultRoute=false when no default route")
	}
}
