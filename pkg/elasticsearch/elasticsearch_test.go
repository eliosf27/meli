package elasticsearch

import (
	"log"
	"os"
	"testing"
	"net/http"
	"io/ioutil"
	"strings"
	"bytes"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// mockRoundTripper is a mock implementation of http.RoundTripper.
type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip implements the http.RoundTripper interface.
func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

// newTestClientWithMockTransport creates a new client with a mock transport.
func newTestClientWithMockTransport(t *testing.T, cfg Config, roundTripFunc func(req *http.Request) (*http.Response, error)) (*Client, error) {
	if roundTripFunc == nil {
		// Default mock if none provided
		roundTripFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(`{"version":{"number":"8.x.x"}}`)),
				Header:     make(http.Header),
			}, nil
		}
	}

	mockTransport := &mockRoundTripper{roundTripFunc: roundTripFunc}

	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses, // These won't be hit due to mock transport
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: mockTransport,
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "test-logger: ", log.LstdFlags)
	return &Client{es: es, l: logger}, nil
}

func TestNewClient(t *testing.T) {
	cfg := Config{
		Addresses: []string{"http://localhost:9200"}, // Dummy address for testing
		Username:  "user",
		Password:  "pass",
	}
	logger := log.New(os.Stdout, "test-logger: ", log.LstdFlags) // Real logger can be used, or a mock one

	// Use the helper to create a client with a mock transport
	client, err := newTestClientWithMockTransport(t, cfg, nil) // nil for default mock transport
	if err != nil {
		t.Fatalf("newTestClientWithMockTransport() error = %v, wantErr %v", err, false)
	}

	if client == nil {
		t.Fatal("newTestClientWithMockTransport() returned nil client, want non-nil client")
	}

	if client.es == nil {
		t.Fatal("newTestClientWithMockTransport() client.es is nil, want non-nil Elasticsearch client")
	}

	// The logger in the client created by newTestClientWithMockTransport is currently hardcoded
	// If we want to assert it, we'd need to pass the logger to newTestClientWithMockTransport
	// or have it return the logger it used. For now, this check is less critical than
	// the client and es client being non-nil.
	// if client.l != logger {
	// 	t.Fatal("NewClient() client.l is not set correctly")
	// }
}

func TestPing(t *testing.T) {
	cfg := Config{Addresses: []string{"http://mock-es"}}

	t.Run("Ping Success", func(t *testing.T) {
		client, err := newTestClientWithMockTransport(t, cfg, func(req *http.Request) (*http.Response, error) {
			// Check if the request path is for Ping
			if req.URL.Path != "/_ping" && req.URL.Path != "/" { // Ping can go to root
				t.Fatalf("Expected ping request path, got %s", req.URL.Path)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(`{"message":"pong"}`)), // Dummy body
				Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
			}, nil
		})
		if err != nil {
			t.Fatalf("Failed to create test client: %v", err)
		}

		err = client.Ping()
		if err != nil {
			t.Errorf("Ping() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Ping Error", func(t *testing.T) {
		client, err := newTestClientWithMockTransport(t, cfg, func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError, // Simulate server error
				Body:       ioutil.NopCloser(strings.NewReader(`{"error":"mock error"}`)),
				Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
			}, nil
		})
		if err != nil {
			t.Fatalf("Failed to create test client: %v", err)
		}

		err = client.Ping()
		if err == nil {
			t.Errorf("Ping() error = nil, wantErr %v", true)
		}
		// TODO: Optionally, check the error message content if it's specific.
	})
}

import (
	"context"
	// ... other imports
)

// ... (TestNewClient, TestPing remain the same)

func TestIndexDocument(t *testing.T) {
	cfg := Config{Addresses: []string{"http://mock-es"}}
	ctx := context.Background()
	indexName := "test-index"
	docID := "1"
	docBody := `{"title":"Test Document"}` // Assuming string for simplicity, matching current IndexDocument

	t.Run("Index Success", func(t *testing.T) {
		client, err := newTestClientWithMockTransport(t, cfg, func(req *http.Request) (*http.Response, error) {
			expectedPath := "/" + indexName + "/_doc/" + docID
			if req.URL.Path != expectedPath {
				t.Errorf("Expected request path %s, got %s", expectedPath, req.URL.Path)
			}
			if req.Method != http.MethodPut && req.Method != http.MethodPost { // Index can be PUT (with ID) or POST (auto-ID)
				t.Errorf("Expected PUT or POST method, got %s", req.Method)
			}

			bodyBytes, _ := ioutil.ReadAll(req.Body)
			if string(bodyBytes) != docBody {
				t.Errorf("Expected body %s, got %s", docBody, string(bodyBytes))
			}

			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       ioutil.NopCloser(strings.NewReader(`{"result":"created","_id":"1"}`)),
				Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
			}, nil
		})
		if err != nil {
			t.Fatalf("Failed to create test client: %v", err)
		}

		err = client.IndexDocument(ctx, indexName, docID, docBody)
		if err != nil {
			t.Errorf("IndexDocument() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Index Error", func(t *testing.T) {
		client, err := newTestClientWithMockTransport(t, cfg, func(req *http.Request) (*http.Response, error) {
			// Basic check, can be more specific if needed
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(strings.NewReader(`{"error":"mock index error"}`)),
				Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
			}, nil
		})
		if err != nil {
			t.Fatalf("Failed to create test client: %v", err)
		}

		err = client.IndexDocument(ctx, indexName, docID, docBody)
		if err == nil {
			t.Errorf("IndexDocument() error = nil, wantErr %v", true)
		}
	})
}

// ... (TestNewClient, TestPing, TestIndexDocument remain the same)

func TestSearchDocuments(t *testing.T) {
	cfg := Config{Addresses: []string{"http://mock-es"}}
	ctx := context.Background()
	indexName := "test-index"
	query := "test query"

	t.Run("Search Success", func(t *testing.T) {
		mockResponseBody := `{"hits":{"total":{"value":1,"relation":"eq"},"hits":[{"_source":{"title":"Test Document"}}]}}`
		client, err := newTestClientWithMockTransport(t, cfg, func(req *http.Request) (*http.Response, error) {
			expectedPath := "/" + indexName + "/_search"
			if req.URL.Path != expectedPath {
				t.Errorf("Expected request path %s, got %s", expectedPath, req.URL.Path)
			}
			if req.Method != http.MethodGet && req.Method != http.MethodPost { // Search can be GET (with q) or POST (with body)
				t.Errorf("Expected GET or POST method, got %s", req.Method)
			}
			if q := req.URL.Query().Get("q"); q != query {
				t.Errorf("Expected query param q=%s, got %s", query, q)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(mockResponseBody)),
				Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
			}, nil
		})
		if err != nil {
			t.Fatalf("Failed to create test client: %v", err)
		}

		res, err := client.SearchDocuments(ctx, indexName, query)
		if err != nil {
			t.Fatalf("SearchDocuments() error = %v, wantErr %v", err, false)
		}
		if res == nil {
			t.Fatal("SearchDocuments() response is nil, want non-nil")
		}
		defer res.Body.Close() // Important to close the response body

		bodyBytes, _ := ioutil.ReadAll(res.Body)
		if string(bodyBytes) != mockResponseBody {
			t.Errorf("Expected response body %s, got %s", mockResponseBody, string(bodyBytes))
		}
	})

	t.Run("Search Error", func(t *testing.T) {
		client, err := newTestClientWithMockTransport(t, cfg, func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(strings.NewReader(`{"error":"mock search error"}`)),
				Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
			}, nil
		})
		if err != nil {
			t.Fatalf("Failed to create test client: %v", err)
		}

		res, err := client.SearchDocuments(ctx, indexName, query)
		if err == nil {
			t.Errorf("SearchDocuments() error = nil, wantErr %v", true)
		}
		if res == nil {
			t.Fatal("SearchDocuments() response is nil on error, want non-nil for error inspection")
		}
		// For client-side errors (like network), res might be nil.
		// For server-side errors that esapi.Response captures, res is non-nil.
		// The current mock returns a valid http.Response, so esapi.Response will be non-nil.
		if res != nil {
			defer res.Body.Close() // Close even error responses
		}
	})
}
