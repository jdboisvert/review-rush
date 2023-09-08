package slack

import (
	"context"
	"testing"
)

type mockService struct{}

func (m *mockService) PostMessage(ctx context.Context, message string) error {
	return nil
}

func TestNewSlackService(t *testing.T) {
	s := NewSlackService("test-token", "test-url")
	if s == nil {
		t.Fatal("expected SlackService, got nil")
	}
}

func TestPostMessage(t *testing.T) {
	s := &service{}

	err := s.PostMessage(context.Background(), "Test message")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
