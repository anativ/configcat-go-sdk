package configcat

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServiceNameInUserAgent(t *testing.T) {
	tests := []struct {
		name              string
		serviceName       string
		pollingMode       PollingMode
		expectedUserAgent string
	}{
		{
			name:              "with service name - AutoPoll",
			serviceName:       "my-service",
			pollingMode:       AutoPoll,
			expectedUserAgent: "ConfigCat-Go/my-service-a-",
		},
		{
			name:              "with service name - Manual",
			serviceName:       "test-app",
			pollingMode:       Manual,
			expectedUserAgent: "ConfigCat-Go/test-app-m-",
		},
		{
			name:              "with service name - Lazy",
			serviceName:       "backend-service",
			pollingMode:       Lazy,
			expectedUserAgent: "ConfigCat-Go/backend-service-l-",
		},
		{
			name:              "without service name - AutoPoll",
			serviceName:       "",
			pollingMode:       AutoPoll,
			expectedUserAgent: "ConfigCat-Go/anativ-a-",
		},
		{
			name:              "without service name - Manual",
			serviceName:       "",
			pollingMode:       Manual,
			expectedUserAgent: "ConfigCat-Go/anativ-m-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedUserAgent string
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedUserAgent = r.Header.Get("X-ConfigCat-UserAgent")
				w.Header().Set("ETag", "test-etag")
				w.WriteHeader(http.StatusOK)
				// Return minimal valid config
				w.Write([]byte(`{"f":{},"p":{"u":"` + server.URL + `","r":0}}`))
			}))
			defer server.Close()

			// Generate a valid SDK key format (22 chars / 22 chars)
			sdkKey := "1234567890123456789012/1234567890123456789012"

			client := NewCustomClient(Config{
				SDKKey:       sdkKey,
				BaseURL:      server.URL,
				ServiceName:  tt.serviceName,
				PollingMode:  tt.pollingMode,
				PollInterval: 1 * time.Hour, // Long interval to prevent auto-polling
			})
			defer client.Close()

			// Trigger a refresh to make an HTTP request
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = client.Refresh(ctx)

			// Give it a moment to complete
			time.Sleep(100 * time.Millisecond)

			if capturedUserAgent == "" {
				t.Fatal("User-Agent was not captured - request may not have been made")
			}

			if !strings.HasPrefix(capturedUserAgent, tt.expectedUserAgent) {
				t.Errorf("Expected User-Agent to start with %q, got %q", tt.expectedUserAgent, capturedUserAgent)
			}

			// Verify it ends with version
			if !strings.Contains(capturedUserAgent, "-"+version) {
				t.Errorf("Expected User-Agent to contain version %q, got %q", version, capturedUserAgent)
			}
		})
	}
}
