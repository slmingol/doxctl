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
	"sort"
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
	serviceHealthCheckWithDeps(conf)
}

// serviceHealthCheckWithDeps allows dependency injection for testing
func serviceHealthCheckWithDeps(config *config) {
	result := svcHealthOutput{
		Timestamp: time.Now(),
		Results:   []svcHealthResult{},
	}

	// Build list of all endpoints to test
	type endpointInfo struct {
		serviceName string
		endpoint    string
		insecure    bool
	}
	var endpoints []endpointInfo
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

		// Determine if we should skip TLS verification (global flag or per-service)
		skipTLS := svcsInsecure || svc.Insecure

		// For each server in the service
		for _, svr := range svc.Svrs {
			// Expand brace expressions - handle nested braces by calling expand multiple times
			expanded := expandBraces(svr)

			for _, expandedSvr := range expanded {
				// Construct service endpoint URL
				endpoint := fmt.Sprintf("https://%s:%d%s", expandedSvr, port, path)
				endpoints = append(endpoints, endpointInfo{serviceName: svc.Svc, endpoint: endpoint, insecure: skipTLS})
			}
		}
	}

	// Test each endpoint with progressive spinner
	err := RunWithSpinnerProgress("Checking service health endpoints", len(endpoints), func(index int) error {
		ep := endpoints[index]
		healthResult := checkServiceEndpoint(ep.serviceName, ep.endpoint, ep.insecure, svcsTimeout)
		result.Results = append(result.Results, healthResult)

		if healthResult.Healthy {
			result.Summary.Healthy++
		} else {
			result.Summary.Failed++
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
func checkServiceEndpoint(serviceName, endpoint string, insecure bool, timeout int) svcHealthResult {
	return checkServiceEndpointWithClient(serviceName, endpoint, insecure, timeout, nil)
}

// checkServiceEndpointWithClient allows dependency injection for testing
func checkServiceEndpointWithClient(serviceName, endpoint string, insecure bool, timeout int, client HTTPClient) svcHealthResult {
	result := svcHealthResult{
		Timestamp: time.Now(),
		Service:   serviceName,
		Endpoint:  endpoint,
		Healthy:   false,
	}

	// Create HTTP client if not provided (for production use)
	if client == nil {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}
		client = &http.Client{
			Timeout:   time.Duration(timeout) * time.Second,
			Transport: transport,
		}
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
	var separators []TableSeparator

	// Group results by service name, then by datacenter
	type datacenterGroup struct {
		datacenter string
		results    []svcHealthResult
	}
	type serviceGroup struct {
		name        string
		datacenters map[string]*datacenterGroup
		dcOrder     []string
	}

	groups := make(map[string]*serviceGroup)
	groupOrder := []string{}

	for _, r := range result.Results {
		if groups[r.Service] == nil {
			groups[r.Service] = &serviceGroup{
				name:        r.Service,
				datacenters: make(map[string]*datacenterGroup),
				dcOrder:     []string{},
			}
			groupOrder = append(groupOrder, r.Service)
		}

		// Extract datacenter from endpoint
		datacenter := extractDatacenterFromEndpoint(r.Endpoint)

		if groups[r.Service].datacenters[datacenter] == nil {
			groups[r.Service].datacenters[datacenter] = &datacenterGroup{
				datacenter: datacenter,
				results:    []svcHealthResult{},
			}
			groups[r.Service].dcOrder = append(groups[r.Service].dcOrder, datacenter)
		}

		groups[r.Service].datacenters[datacenter].results = append(groups[r.Service].datacenters[datacenter].results, r)
	}

	// Sort services alphabetically
	sort.Strings(groupOrder)

	// Build rows with separators between service and datacenter groups
	rowCount := 0
	for svcIdx, serviceName := range groupOrder {
		group := groups[serviceName]

		// Sort datacenters alphabetically
		sort.Strings(group.dcOrder)

		// Count total items in this service
		totalItemsInService := 0
		for _, dc := range group.dcOrder {
			totalItemsInService += len(group.datacenters[dc].results)
		}

		// If service has <= 5 items, no separators within service
		// If > 5 items, group into chunks of 4-5 at datacenter boundaries
		itemsInCurrentChunk := 0

		for dcIdx, dc := range group.dcOrder {
			dcGroup := group.datacenters[dc]

			// Sort results within each datacenter alphabetically by endpoint
			sort.Slice(dcGroup.results, func(i, j int) bool {
				return dcGroup.results[i].Endpoint < dcGroup.results[j].Endpoint
			})

			for _, r := range dcGroup.results {
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
				rowCount++
				itemsInCurrentChunk++
			}

			// Add light separator if:
			// 1. This service has > 5 items total
			// 2. We've accumulated 4-5 items in current chunk
			// 3. This isn't the last datacenter in the service
			if totalItemsInService > 5 && dcIdx < len(group.dcOrder)-1 {
				if itemsInCurrentChunk >= 4 {
					separators = append(separators, TableSeparator{
						RowIndex: rowCount - 1,
						Type:     LightSeparator,
					})
					itemsInCurrentChunk = 0
				}
			}
		}

		// Add heavy separator after each service group (except last)
		if svcIdx < len(groupOrder)-1 {
			separators = append(separators, TableSeparator{
				RowIndex: rowCount - 1,
				Type:     HeavySeparator,
			})
		}
	}

	fmt.Print(createStyledTableWithTypedSeparators(headers, rows, "Service Health Checks (HTTPS)", separators))

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

// extractDatacenterFromEndpoint extracts the datacenter identifier from an endpoint URL
func extractDatacenterFromEndpoint(endpoint string) string {
	// Remove protocol
	host := endpoint
	if strings.HasPrefix(host, "https://") {
		host = host[8:]
	} else if strings.HasPrefix(host, "http://") {
		host = host[7:]
	}

	// Remove path and port
	if idx := strings.Index(host, "/"); idx != -1 {
		host = host[:idx]
	}
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Check if it's an IP address
	if strings.Count(host, ".") == 3 {
		parts := strings.Split(host, ".")
		if len(parts) == 4 {
			isIP := true
			for _, p := range parts {
				if len(p) == 0 || len(p) > 3 {
					isIP = false
					break
				}
			}
			if isIP {
				return "ip"
			}
		}
	}

	parts := strings.Split(host, ".")
	if len(parts) >= 3 {
		// For patterns like:
		// - api.app1.lab1.ocp.bandwidth.com -> lab1
		// - es-master-01d.lab1.bwnet.us -> lab1
		// - idm-01a.lab1.bandwidthclec.local -> lab1
		// - idm-01a.bru1.bwnet.us -> bru1

		if strings.HasPrefix(host, "api.app1.") {
			// api.app1.lab1.ocp... -> lab1
			return parts[2]
		} else if strings.HasPrefix(host, "es-master-") {
			// es-master-01d.lab1.bwnet.us -> lab1
			return parts[1]
		} else if strings.HasPrefix(host, "idm-") {
			// idm-01a.lab1.bandwidthclec.local -> lab1
			return parts[1]
		}
	}

	return "unknown"
}
