/*
Package cmd - Testable interfaces for dependency injection

Copyright © 2021 Sam Mingolelli <github@lamolabs.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/go-ping/ping"
	gobrex "github.com/kujtimiihoxha/go-brace-expansion"
)

// DNSResolver interface for DNS lookup operations
type DNSResolver interface {
	LookupHost(ctx context.Context, host string) ([]string, error)
}

// Pinger interface for ping operations
type Pinger interface {
	SetTimeout(duration time.Duration)
	SetCount(count int)
	Run() error
	Statistics() *ping.Statistics
}

// CommandExecutor interface for running system commands
type CommandExecutor interface {
	Execute(name string, args ...string) ([]byte, error)
}

// FileReader interface for reading files
type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

// BraceExpander interface for brace expansion
type BraceExpander interface {
	Expand(pattern string) []string
}

// ========== Real Implementations ==========

// realDNSResolver wraps the standard net.Resolver
type realDNSResolver struct{}

func (r *realDNSResolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	resolver := &net.Resolver{}
	return resolver.LookupHost(ctx, host)
}

// realPinger wraps the go-ping Pinger
type realPinger struct {
	inner *ping.Pinger
}

func (p *realPinger) SetTimeout(duration time.Duration) {
	p.inner.Timeout = duration
}

func (p *realPinger) SetCount(count int) {
	p.inner.Count = count
}

func (p *realPinger) Run() error {
	return p.inner.Run()
}

func (p *realPinger) Statistics() *ping.Statistics {
	return p.inner.Statistics()
}

// realCommandExecutor wraps exec.Command
type realCommandExecutor struct{}

func (e *realCommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

// realFileReader wraps os.ReadFile
type realFileReader struct{}

func (r *realFileReader) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// realBraceExpander wraps the gobrex package
type realBraceExpander struct{}

func (e *realBraceExpander) Expand(pattern string) []string {
	return gobrex.Expand(pattern)
}

// ========== Factory Functions (can be overridden in tests) ==========

var (
	// NewDNSResolver creates a new DNS resolver (can be mocked in tests)
	NewDNSResolver = func() DNSResolver {
		return &realDNSResolver{}
	}

	// NewPinger creates a new pinger (can be mocked in tests)
	NewPinger = func(host string) (Pinger, error) {
		p, err := ping.NewPinger(host)
		if err != nil {
			return nil, err
		}
		return &realPinger{inner: p}, nil
	}

	// NewCommandExecutor creates a new command executor (can be mocked in tests)
	NewCommandExecutor = func() CommandExecutor {
		return &realCommandExecutor{}
	}

	// NewFileReader creates a new file reader (can be mocked in tests)
	NewFileReader = func() FileReader {
		return &realFileReader{}
	}

	// NewBraceExpander creates a new brace expander (can be mocked in tests)
	NewBraceExpander = func() BraceExpander {
		return &realBraceExpander{}
	}
)
