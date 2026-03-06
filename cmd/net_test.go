package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/go-ping/ping"
)

type mockNetPinger struct {
	stats *ping.Statistics
	err   error
}

func (m *mockNetPinger) Run() error {
	return m.err
}

func (m *mockNetPinger) Statistics() *ping.Statistics {
	return m.stats
}

func (m *mockNetPinger) SetTimeout(timeout time.Duration) {}

func (m *mockNetPinger) SetCount(count int) {}

func TestNetPerformanceCheckWithDeps_Success(t *testing.T) {
	config := &config{
		Sites:      []string{"site1", "site2"},
		DomainName: "example.com",
	}

	pingerFactory := func(addr string) (Pinger, error) {
		return &mockNetPinger{
			stats: &ping.Statistics{
				PacketsRecv: 10,
				PacketsSent: 10,
				PacketLoss:  0.0,
				AvgRtt:      30 * time.Millisecond,
				MinRtt:      25 * time.Millisecond,
				MaxRtt:      35 * time.Millisecond,
				StdDevRtt:   3 * time.Millisecond,
			},
			err: nil,
		}, nil
	}

	// Should not panic
	netPerformanceCheckWithDeps(config, 50.0, 10, pingerFactory)
}

func TestNetPerformanceCheckWithDeps_HighLatency(t *testing.T) {
	config := &config{
		Sites:      []string{"site1"},
		DomainName: "example.com",
	}

	pingerFactory := func(addr string) (Pinger, error) {
		return &mockNetPinger{
			stats: &ping.Statistics{
				PacketsRecv: 10,
				PacketsSent: 10,
				PacketLoss:  0.0,
				AvgRtt:      100 * time.Millisecond, // Exceeds 50ms SLO
				MinRtt:      90 * time.Millisecond,
				MaxRtt:      110 * time.Millisecond,
				StdDevRtt:   5 * time.Millisecond,
			},
			err: nil,
		}, nil
	}

	// Should not panic, but SLO should fail
	netPerformanceCheckWithDeps(config, 50.0, 10, pingerFactory)
}

func TestNetPerformanceCheckWithDeps_PacketLoss(t *testing.T) {
	config := &config{
		Sites:      []string{"site1"},
		DomainName: "example.com",
	}

	pingerFactory := func(addr string) (Pinger, error) {
		return &mockNetPinger{
			stats: &ping.Statistics{
				PacketsRecv: 7,
				PacketsSent: 10,
				PacketLoss:  30.0, // High packet loss
				AvgRtt:      30 * time.Millisecond,
				MinRtt:      25 * time.Millisecond,
				MaxRtt:      35 * time.Millisecond,
				StdDevRtt:   3 * time.Millisecond,
			},
			err: nil,
		}, nil
	}

	// Should not panic, packet loss should cause SLO fail
	netPerformanceCheckWithDeps(config, 50.0, 10, pingerFactory)
}

func TestNetPerformanceCheckWithDeps_PingError(t *testing.T) {
	config := &config{
		Sites:      []string{"site1"},
		DomainName: "example.com",
	}

	pingerFactory := func(addr string) (Pinger, error) {
		return &mockNetPinger{
			stats: &ping.Statistics{
				PacketsRecv: 0,
				PacketsSent: 10,
				PacketLoss:  100.0,
			},
			err: errors.New("network unreachable"),
		}, nil
	}

	// Should not panic
	netPerformanceCheckWithDeps(config, 50.0, 10, pingerFactory)
}

func TestNetPerformanceCheckWithDeps_NoTargets(t *testing.T) {
	config := &config{
		Sites:      []string{},
		DomainName: "example.com",
	}

	pingerFactory := func(addr string) (Pinger, error) {
		return nil, errors.New("should not be called")
	}

	// Should handle gracefully
	netPerformanceCheckWithDeps(config, 50.0, 10, pingerFactory)
}

func TestNetPerformanceCheckWithDeps_PingerFactoryError(t *testing.T) {
	config := &config{
		Sites:      []string{"site1"},
		DomainName: "example.com",
	}

	pingerFactory := func(addr string) (Pinger, error) {
		return nil, errors.New("failed to create pinger")
	}

	// Should skip target gracefully
	netPerformanceCheckWithDeps(config, 50.0, 10, pingerFactory)
}

func TestNetPerformanceCheckWithDeps_MultipleTargets(t *testing.T) {
	config := &config{
		Sites:      []string{"site1", "site2", "site3"},
		DomainName: "example.com",
	}

	callCount := 0
	pingerFactory := func(addr string) (Pinger, error) {
		callCount++
		latency := time.Duration(callCount*20) * time.Millisecond

		return &mockNetPinger{
			stats: &ping.Statistics{
				PacketsRecv: 10,
				PacketsSent: 10,
				PacketLoss:  0.0,
				AvgRtt:      latency,
				MinRtt:      latency - 5*time.Millisecond,
				MaxRtt:      latency + 5*time.Millisecond,
				StdDevRtt:   2 * time.Millisecond,
			},
			err: nil,
		}, nil
	}

	// Should process all targets
	netPerformanceCheckWithDeps(config, 50.0, 10, pingerFactory)

	if callCount != 3 {
		t.Errorf("Expected 3 pinger calls, got %d", callCount)
	}
}

func TestNetPerformanceCheckWithDeps_ZeroPackets(t *testing.T) {
	config := &config{
		Sites:      []string{"site1"},
		DomainName: "example.com",
	}

	pingerFactory := func(addr string) (Pinger, error) {
		return &mockNetPinger{
			stats: &ping.Statistics{
				PacketsRecv: 0,
				PacketsSent: 10,
				PacketLoss:  100.0,
				AvgRtt:      0,
				MinRtt:      0,
				MaxRtt:      0,
				StdDevRtt:   0,
			},
			err: nil,
		}, nil
	}

	// Should handle zero received packets
	netPerformanceCheckWithDeps(config, 50.0, 10, pingerFactory)
}
