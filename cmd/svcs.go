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
	"time"

	"doxctl/internal/output"

	"github.com/jedib0t/go-pretty/v6/table"
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

	// Iterate through configured services
	for _, svc := range config.Svcs {
		// For each server in the service
		for _, svr := range svc.Svrs {
			// Construct service endpoint URL
			// Assuming OpenShift API endpoint pattern
			endpoint := fmt.Sprintf("https://%s:6443/healthz", svr)

			healthResult := checkServiceEndpoint(svc.Svc, endpoint, client)
			result.Results = append(result.Results, healthResult)

			if healthResult.Healthy {
				result.Summary.Healthy++
			} else {
				result.Summary.Failed++
			}
		}
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
	fmt.Println()

	t := table.NewWriter()
	t.SetTitle("Service Health Checks")
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Service", "Endpoint", "Status Code", "Response (ms)", "Status", "Error"})

	for _, r := range result.Results {
		status := "✓ Healthy"
		if !r.Healthy {
			status = "✗ Failed"
		}

		statusCodeStr := "-"
		if r.StatusCode > 0 {
			statusCodeStr = fmt.Sprintf("%d", r.StatusCode)
		}

		responseTimeStr := "-"
		if r.ResponseTimeMs > 0 {
			responseTimeStr = fmt.Sprintf("%.2f", r.ResponseTimeMs)
		}

		errorStr := ""
		if r.Error != "" {
			// Trim long error messages
			if len(r.Error) > 40 {
				errorStr = r.Error[:37] + "..."
			} else {
				errorStr = r.Error
			}
		}

		t.AppendRow([]interface{}{
			r.Service,
			r.Endpoint,
			statusCodeStr,
			responseTimeStr,
			status,
			errorStr,
		})
	}

	t.Render()

	// Print summary
	fmt.Printf("\nSummary: %d/%d services healthy (%.1f%% availability)\n",
		result.Summary.Healthy,
		result.Summary.Total,
		float64(result.Summary.Healthy)/float64(result.Summary.Total)*100)
	fmt.Println()
}
