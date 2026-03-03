/*
Package cmd - Service-level endpoint health checks

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
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"doxctl/internal/output"

	"github.com/spf13/cobra"
)

// Service health check result
type svcHealthResult struct {
	Timestamp      time.Time `json:"timestamp" yaml:"timestamp"`
	Service        string    `json:"service" yaml:"service"`
	Endpoint       string    `json:"endpoint" yaml:"endpoint"`
	ResponseTimeMs float64   `json:"responseTimeMs" yaml:"responseTimeMs"`
	StatusCode     int       `json:"statusCode" yaml:"statusCode"`
	Healthy        bool      `json:"healthy" yaml:"healthy"`
	Error          string    `json:"error,omitempty" yaml:"error,omitempty"`
}

type svcHealthOutput struct {
	Timestamp time.Time         `json:"timestamp" yaml:"timestamp"`
	Results   []svcHealthResult `json:"results" yaml:"results"`
	Summary   struct {
		Total   int `json:"total" yaml:"total"`
		Healthy int `json:"healthy" yaml:"healthy"`
		Failed  int `json:"failed" yaml:"failed"`
	} `json:"summary" yaml:"summary"`
}

var (
	svcsHealthChk bool
	svcsTimeout   int
	svcsInsecure  bool
)

// HTTPClient interface for dependency injection
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// svcsCmd represents the svcs command
var svcsCmd = &cobra.Command{
	Use:   "svcs",
	Short: "Service-level health checks for multi-datacenter endpoints",
	Long: `Check the health and availability of services across multiple datacenters.

This command performs HTTP/HTTPS health checks on service endpoints and measures:
  - Response time
  - HTTP status codes
  - Service availability
  - Multi-datacenter service health

Examples:
  # Check health of all configured services
  doxctl svcs --health

  # Set custom timeout (default: 5 seconds)
  doxctl svcs --health --timeout 10

  # Skip TLS verification for self-signed certificates
  doxctl svcs --health --insecure`,
	Run: func(cmd *cobra.Command, args []string) {
		if allChk || svcsHealthChk {
			serviceHealthCheck()
		} else {
			_ = cmd.Usage()
			fmt.Printf("\n")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(svcsCmd)

	svcsCmd.Flags().BoolVarP(&svcsHealthChk, "health", "H", false, "Run service health checks")
	svcsCmd.Flags().IntVarP(&svcsTimeout, "timeout", "t", 5, "HTTP request timeout in seconds")
	svcsCmd.Flags().BoolVarP(&svcsInsecure, "insecure", "k", false, "Skip TLS certificate verification")
}

// serviceHealthCheck checks the health of configured services
func serviceHealthCheck() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: svcsInsecure},
	}
	client := &http.Client{
		Timeout:   time.Duration(svcsTimeout) * time.Second,
		Transport: transport,
	}

	serviceHealthCheckWithDeps(conf, client)
}

// serviceHealthCheckWithDeps allows dependency injection for testing
func serviceHealthCheckWithDeps(config *config, client HTTPClient) {
	result := svcHealthOutput{
		Timestamp: time.Now(),
		Results:   []svcHealthResult{},
	}

	// Wrap health checks in a spinner
	err := RunWithSpinner("Checking service health endpoints", func() error {
		// Iterate through configured services
		for _, svc := range config.Svcs {
		// Default port to 6443 if not specified
		port := svc.Port
		if port == 0 {
			port = 6443
		}

		// Default path to /healthz if not specified
		path := svc.Path
		if path == "" {
			path = "/healthz"
		}

		// For each server in the service
		for _, svr := range svc.Svrs {
			// Expand brace expressions - handle nested braces by calling expand multiple times
			expanded := expandBraces(svr)

			for _, expandedSvr := range expanded {
				// Construct service endpoint URL
				endpoint := fmt.Sprintf("https://%s:%d%s", expandedSvr, port, path)

				healthResult := checkServiceEndpoint(svc.Svc, endpoint, client)
				result.Results = append(result.Results, healthResult)

				if healthResult.Healthy {
					result.Summary.Healthy++
				} else {
					result.Summary.Failed++
				}
			}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during health checks: %v\n", err)
		os.Exit(1)
	}

	result.Summary.Total = len(result.Results)

	// Output results
	switch outputFormat {
	case "json":
		output.PrintJSON(result)
	case "yaml":
		output.PrintYAML(result)
	default:
		printSvcHealthTable(result)
	}
}

// expandBraces recursively expands nested brace expressions
func expandBraces(s string) []string {
	// Find first opening brace
	start := strings.Index(s, "{")
	if start == -1 {
		// No braces, return as-is
		return []string{s}
	}

	// Find matching closing brace
	depth := 0
	end := -1
	for i := start; i < len(s); i++ {
		if s[i] == '{' {
			depth++
		} else if s[i] == '}' {
			depth--
			if depth == 0 {
				end = i
				break
			}
		}
	}

	if end == -1 {
		// No matching brace, return as-is
		return []string{s}
	}

	// Extract the options between braces
	options := strings.Split(s[start+1:end], ",")
	prefix := s[:start]
	suffix := s[end+1:]

	// Expand this level
	var results []string
	for _, opt := range options {
		results = append(results, prefix+opt+suffix)
	}

	// Recursively expand each result (for nested braces)
	var finalResults []string
	for _, r := range results {
		finalResults = append(finalResults, expandBraces(r)...)
	}

	return finalResults
}

// checkServiceEndpoint performs a single health check against an endpoint
func checkServiceEndpoint(serviceName, endpoint string, client HTTPClient) svcHealthResult {
	result := svcHealthResult{
		Timestamp: time.Now(),
		Service:   serviceName,
		Endpoint:  endpoint,
		Healthy:   false,
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	result.ResponseTimeMs = float64(elapsed.Microseconds()) / 1000.0

	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// Consider 2xx and 3xx status codes as healthy
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		result.Healthy = true
	}

	return result
}

func printSvcHealthTable(result svcHealthOutput) {
	// Table headers - compact version
	headers := []string{"Service", "Endpoint", "Status", "Response (ms)"}
	if verboseChk {
		headers = append(headers, "Error")
	}

	var rows [][]string
	for _, r := range result.Results {
		// Status icon
		status := "✓"
		if !r.Healthy {
			status = "✗"
		}

		// Response time
		responseTimeStr := "-"
		if r.ResponseTimeMs > 0 {
			responseTimeStr = fmt.Sprintf("%.2f", r.ResponseTimeMs)
		}

		// Extract hostname:port from endpoint URL (remove protocol and path)
		host := r.Endpoint
		// Remove https:// or http:// prefix
		if strings.HasPrefix(host, "https://") {
			host = host[8:]
		} else if strings.HasPrefix(host, "http://") {
			host = host[7:]
		}
		// Remove path (keep hostname:port)
		if idx := strings.Index(host, "/"); idx != -1 {
			host = host[:idx]
		}

		row := []string{
			r.Service,
			host,
			status,
			responseTimeStr,
		}

		// Add error column only in verbose mode
		if verboseChk {
			errorStr := ""
			if r.Error != "" {
				// In verbose mode, show more of the error (truncate at 120 chars)
				if len(r.Error) > 120 {
					errorStr = r.Error[:117] + "..."
				} else {
					errorStr = r.Error
				}
			}
			row = append(row, errorStr)
		}

		rows = append(rows, row)
	}

	fmt.Print(createStyledTable(headers, rows, "Service Health Checks (HTTPS)"))

	// Print summary
	availability := 0.0
	if result.Summary.Total > 0 {
		availability = float64(result.Summary.Healthy) / float64(result.Summary.Total) * 100
	}
	fmt.Printf("\nSummary: %d/%d services healthy (%.1f%% availability)\n\n",
		result.Summary.Healthy,
		result.Summary.Total,
		availability)
}
