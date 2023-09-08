package slack

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func newMockClient(response *http.Response, err error) *http.Client {
	return &http.Client{
		Transport: &mockTransport{
			response: response,
			err:      err,
		},
	}
}

func TestNewSlackService(t *testing.T) {
	token := "testToken"
	url := "testUrl"
	slackService := NewSlackService(token, url)

	impl, ok := slackService.(*service)
	if !ok {
		t.Fatalf("Expected type '*service', got '%T'", slackService)
	}

	if impl.token != token {
		t.Errorf("Expected token '%s', got '%s'", token, impl.token)
	}

	if impl.url != url {
		t.Errorf("Expected URL '%s', got '%s'", url, impl.url)
	}

	if impl.client == nil {
		t.Error("Expected client to be initialized, got nil")
	}
}

func TestPostMessage(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		clientErr  error
		expectErr  bool
	}{
		{
			name:       "successful post",
			response:   `{"ok": true}`,
			statusCode: 200,
			expectErr:  false,
		},
		{
			name:       "failed post with Slack error",
			response:   `{"ok": false, "error": "invalid_token"}`,
			statusCode: 200,
			expectErr:  true,
		},
		{
			name:       "failed post with non-200 status",
			response:   "",
			statusCode: 400,
			expectErr:  true,
		},
		{
			name:      "client error",
			response:  "",
			clientErr: fmt.Errorf("some error"),
			expectErr: true,
		},
	}

	token := "test_token"
	url := "test_url"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(strings.NewReader(tt.response)),
			}
			client := newMockClient(resp, tt.clientErr)
			svc := &service{
				client: client,
				token:  token,
				url:    url,
			}

			err := svc.PostMessage(context.TODO(), "some_channel", "hello")

			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
