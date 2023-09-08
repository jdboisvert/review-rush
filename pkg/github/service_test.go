package github

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type MockTransport struct {
	Resp *http.Response
	Err  error
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Resp, m.Err
}

func TestNewGitHubService(t *testing.T) {
	token := "testToken"
	org := "testOrg"

	serviceInstance := NewGitHubService(token, org)

	s, ok := serviceInstance.(*service)
	if !ok {
		t.Fatalf("Expected type *service; got %T", serviceInstance)
	}

	if s.token != token {
		t.Errorf("Expected token to be %s; got %s", token, s.token)
	}

	if s.org != org {
		t.Errorf("Expected org to be %s; got %s", org, s.org)
	}

	if reflect.TypeOf(s.client).String() != "*http.Client" {
		t.Errorf("Expected client to be of type *http.Client; got %T", s.client)
	}
}

func TestGetRepos(t *testing.T) {
	responseBody := `[{"id": 1, "name": "repo1"}, {"id": 2, "name": "repo2"}]`
	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(responseBody))),
	}

	client := &http.Client{
		Transport: &MockTransport{Resp: mockResp},
	}

	service := &service{
		client: client,
		token:  "sample_token",
		org:    "sample_org",
	}

	repos, err := service.GetRepos(context.Background())
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	if len(repos) != 2 {
		t.Fatalf("expected 2 repos but got %d", len(repos))
	}
}

func TestGetPullRequests(t *testing.T) {
	responseBody := `[{"id": 1, "title": "PR1"}, {"id": 2, "title": "PR2"}]`
	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(responseBody))),
	}

	client := &http.Client{
		Transport: &MockTransport{Resp: mockResp},
	}

	service := &service{
		client: client,
		token:  "sample_token",
		org:    "sample_org",
	}

	prs, err := service.GetPullRequests(context.Background(), "sample_repo")
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	if len(prs) != 2 {
		t.Fatalf("expected 2 pull requests but got %d", len(prs))
	}
}

func TestGetReviews(t *testing.T) {
	responseBody := `[{"id": 1, "user": {"login": "user1"}}, {"id": 2, "user": {"login": "user2"}}]`
	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(responseBody))),
	}

	client := &http.Client{
		Transport: &MockTransport{Resp: mockResp},
	}

	service := &service{
		client: client,
		token:  "sample_token",
		org:    "sample_org",
	}

	reviews, err := service.GetReviews(context.Background(), "https://api.github.com/repos/sample_org/sample_repo/pulls/1")
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	if len(reviews) != 2 {
		t.Fatalf("expected 2 reviews but got %d", len(reviews))
	}
}
