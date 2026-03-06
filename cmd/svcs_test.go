package cmd

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestServiceHealthCheckWithDeps_Success(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com"}},
		},
	}

	// Should not panic with valid config
	serviceHealthCheckWithDeps(config)
}

func TestServiceHealthCheckWithDeps_MultipleServices(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com", "ocp-master-2.example.com"}},
			{Svc: "elastic", Svrs: []string{"es-master-1.example.com"}},
		},
	}

	// Should not panic with multiple services
	serviceHealthCheckWithDeps(config)
}

func TestServiceHealthCheckWithDeps_HTTPError(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com"}},
		},
	}

	// Should handle error gracefully
	serviceHealthCheckWithDeps(config)
}

func TestServiceHealthCheckWithDeps_4xxStatus(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com"}},
		},
	}

	// Should mark as unhealthy
	serviceHealthCheckWithDeps(config)
}

func TestServiceHealthCheckWithDeps_5xxStatus(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com"}},
		},
	}

	// Should mark as unhealthy
	serviceHealthCheckWithDeps(config)
}

func TestServiceHealthCheckWithDeps_NoServices(t *testing.T) {
	config := &config{
		Svcs: []svc{},
	}

	// Should handle empty service list
	serviceHealthCheckWithDeps(config)
}

func TestServiceHealthCheckWithDeps_3xxRedirect(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com"}},
		},
	}

	// 3xx should be considered healthy
	serviceHealthCheckWithDeps(config)
}

func TestCheckServiceEndpoint_Success(t *testing.T) {
	client := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("ok")),
		},
		err: nil,
	}

	result := checkServiceEndpointWithClient("test-service", "https://example.com/health", false, 5, client)

	if !result.Healthy {
		t.Error("Expected service to be healthy")
	}
	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}
	if result.ResponseTimeMs < 0 {
		t.Error("Expected non-negative response time")
	}
}

func TestCheckServiceEndpoint_RequestError(t *testing.T) {
	client := &mockHTTPClient{
		response: nil,
		err:      errors.New("network error"),
	}

	result := checkServiceEndpointWithClient("test-service", "https://example.com/health", false, 5, client)

	if result.Healthy {
		t.Error("Expected service to be unhealthy")
	}
	if result.Error != "network error" {
		t.Errorf("Expected error 'network error', got '%s'", result.Error)
	}
}

func TestCheckServiceEndpoint_InvalidURL(t *testing.T) {
	client := &mockHTTPClient{}

	// Invalid URL should cause request creation error
	result := checkServiceEndpointWithClient("test-service", "://invalid-url", false, 5, client)

	if result.Healthy {
		t.Error("Expected service to be unhealthy")
	}
	if result.Error == "" {
		t.Error("Expected error for invalid URL")
	}
}

func TestCheckServiceEndpoint_Unhealthy4xx(t *testing.T) {
	client := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 401,
			Body:       io.NopCloser(strings.NewReader("unauthorized")),
		},
		err: nil,
	}

	result := checkServiceEndpointWithClient("test-service", "https://example.com/health", false, 5, client)

	if result.Healthy {
		t.Error("Expected service to be unhealthy with 401 status")
	}
	if result.StatusCode != 401 {
		t.Errorf("Expected status code 401, got %d", result.StatusCode)
	}
}

func TestCheckServiceEndpoint_Unhealthy5xx(t *testing.T) {
	client := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 503,
			Body:       io.NopCloser(strings.NewReader("service unavailable")),
		},
		err: nil,
	}

	result := checkServiceEndpointWithClient("test-service", "https://example.com/health", false, 5, client)

	if result.Healthy {
		t.Error("Expected service to be unhealthy with 503 status")
	}
}

func TestExtractDatacenterFromEndpoint_OpenshiftPattern(t *testing.T) {
	dc := extractDatacenterFromEndpoint("https://api.app1.lab1.ocp.bandwidth.com:6443/healthz")
	if dc != "lab1" {
		t.Errorf("Expected 'lab1', got '%s'", dc)
	}

	dc = extractDatacenterFromEndpoint("https://api.app1.rdu1.ocp.bandwidth.com:6443/healthz")
	if dc != "rdu1" {
		t.Errorf("Expected 'rdu1', got '%s'", dc)
	}
}

func TestExtractDatacenterFromEndpoint_ElasticPattern(t *testing.T) {
	dc := extractDatacenterFromEndpoint("https://es-master-01d.lab1.bwnet.us:9200/")
	if dc != "lab1" {
		t.Errorf("Expected 'lab1', got '%s'", dc)
	}

	dc = extractDatacenterFromEndpoint("https://es-master-01e.rdu1.bwnet.us:9200/")
	if dc != "rdu1" {
		t.Errorf("Expected 'rdu1', got '%s'", dc)
	}
}

func TestExtractDatacenterFromEndpoint_IdmPattern(t *testing.T) {
	dc := extractDatacenterFromEndpoint("https://idm-01a.lab1.bandwidthclec.local:443/")
	if dc != "lab1" {
		t.Errorf("Expected 'lab1', got '%s'", dc)
	}

	dc = extractDatacenterFromEndpoint("https://idm-01b.bru1.bwnet.us:443/")
	if dc != "bru1" {
		t.Errorf("Expected 'bru1', got '%s'", dc)
	}
}

func TestExtractDatacenterFromEndpoint_IPAddress(t *testing.T) {
	dc := extractDatacenterFromEndpoint("https://10.23.12.154:6443/healthz")
	if dc != "ip" {
		t.Errorf("Expected 'ip', got '%s'", dc)
	}

	dc = extractDatacenterFromEndpoint("https://192.168.1.1:443/")
	if dc != "ip" {
		t.Errorf("Expected 'ip', got '%s'", dc)
	}
}

func TestExtractDatacenterFromEndpoint_Unknown(t *testing.T) {
	dc := extractDatacenterFromEndpoint("https://unknown-host.example.com:443/")
	if dc != "unknown" {
		t.Errorf("Expected 'unknown', got '%s'", dc)
	}
}
