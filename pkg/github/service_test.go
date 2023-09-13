package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type MockTransport struct {
	Resp     *http.Response
	Err      error
	RespFunc func(req *http.Request) *http.Response
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.RespFunc != nil {
		return m.RespFunc(req), m.Err
	}
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
	firstPage := `[{"id": 1, "name": "repo1"}, {"id": 2, "name": "repo2"}]`
	secondPage := `[]` // no more repos on the next page

	callCount := 0
	transport := &MockTransport{
		RespFunc: func(req *http.Request) *http.Response {
			callCount++
			if callCount == 1 {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader([]byte(firstPage))),
				}
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(secondPage))),
			}
		},
	}

	client := &http.Client{
		Transport: transport,
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

	if callCount != 2 {
		t.Fatalf("expected 2 API calls but got %d", callCount)
	}
}

func TestGetReposEmpty(t *testing.T) {
	firstPage := `[]`

	transport := &MockTransport{
		RespFunc: func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(firstPage))),
			}
		},
	}

	client := &http.Client{
		Transport: transport,
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

	if len(repos) != 0 {
		t.Fatalf("expected 0 repos but got %d", len(repos))
	}
}

func TestGetPullRequestsUpdatedToday(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	responseBody := fmt.Sprintf(`[{"id": 1, "title": "PR1", "updated_at": "%sT12:30:00Z"}, {"id": 2, "title": "PR2", "updated_at": "%sT12:30:00Z"}]`, today, today)
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

func TestGetReviewsDoneToday(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	responseBody := fmt.Sprintf(`[{"id": 1, "user": {"login": "user1"}, "submitted_at": "%sT12:30:00Z"}, {"id": 2, "user": {"login": "user2"}, "submitted_at": "%sT13:45:00Z"}]`, today, today)
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

func TestGetReviewsDoneNotToday(t *testing.T) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	responseBody := fmt.Sprintf(`[{"id": 1, "user": {"login": "user1"}, "submitted_at": "%sT12:30:00Z"}, {"id": 2, "user": {"login": "user2"}, "submitted_at": "%sT13:45:00Z"}]`, yesterday, yesterday)
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

	if len(reviews) != 0 {
		t.Fatalf("expected 0 reviews but got %d", len(reviews))
	}
}
