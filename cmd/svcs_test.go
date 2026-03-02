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

	client := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("ok")),
		},
		err: nil,
	}

	// Should not panic
	serviceHealthCheckWithDeps(config, client)
}

func TestServiceHealthCheckWithDeps_MultipleServices(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com", "ocp-master-2.example.com"}},
			{Svc: "elastic", Svrs: []string{"es-master-1.example.com"}},
		},
	}

	callCount := 0
	client := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("ok")),
		},
		err: nil,
	}

	// Mock to track calls
	origClient := client
	trackingClient := &trackingHTTPClient{
		client:    origClient,
		callCount: &callCount,
	}

	serviceHealthCheckWithDeps(config, trackingClient)

	// Should make 3 calls (2 openshift + 1 elastic)
	if callCount != 3 {
		t.Errorf("Expected 3 HTTP calls, got %d", callCount)
	}
}

type trackingHTTPClient struct {
	client    HTTPClient
	callCount *int
}

func (t *trackingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	*t.callCount++
	return t.client.Do(req)
}

func TestServiceHealthCheckWithDeps_HTTPError(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com"}},
		},
	}

	client := &mockHTTPClient{
		response: nil,
		err:      errors.New("connection refused"),
	}

	// Should handle error gracefully
	serviceHealthCheckWithDeps(config, client)
}

func TestServiceHealthCheckWithDeps_4xxStatus(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com"}},
		},
	}

	client := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(strings.NewReader("not found")),
		},
		err: nil,
	}

	// Should mark as unhealthy
	serviceHealthCheckWithDeps(config, client)
}

func TestServiceHealthCheckWithDeps_5xxStatus(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com"}},
		},
	}

	client := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(strings.NewReader("internal server error")),
		},
		err: nil,
	}

	// Should mark as unhealthy
	serviceHealthCheckWithDeps(config, client)
}

func TestServiceHealthCheckWithDeps_NoServices(t *testing.T) {
	config := &config{
		Svcs: []svc{},
	}

	client := &mockHTTPClient{
		response: nil,
		err:      errors.New("should not be called"),
	}

	// Should handle empty service list
	serviceHealthCheckWithDeps(config, client)
}

func TestServiceHealthCheckWithDeps_3xxRedirect(t *testing.T) {
	config := &config{
		Svcs: []svc{
			{Svc: "openshift", Svrs: []string{"ocp-master-1.example.com"}},
		},
	}

	client := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 301,
			Body:       io.NopCloser(strings.NewReader("moved permanently")),
		},
		err: nil,
	}

	// 3xx should be considered healthy
	serviceHealthCheckWithDeps(config, client)
}

func TestCheckServiceEndpoint_Success(t *testing.T) {
	client := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("ok")),
		},
		err: nil,
	}

	result := checkServiceEndpoint("test-service", "https://example.com/health", client)

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

	result := checkServiceEndpoint("test-service", "https://example.com/health", client)

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
	result := checkServiceEndpoint("test-service", "://invalid-url", client)

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

	result := checkServiceEndpoint("test-service", "https://example.com/health", client)

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

	result := checkServiceEndpoint("test-service", "https://example.com/health", client)

	if result.Healthy {
		t.Error("Expected service to be unhealthy with 503 status")
	}
}
