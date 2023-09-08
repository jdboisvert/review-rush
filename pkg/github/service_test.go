package github

import (
	"context"
	"reflect"
	"testing"
)

func TestGetRepos(t *testing.T) {
	s := NewGitHubService("dummyToken", "dummyOrg")

	repos, err := s.GetRepos(context.Background())
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	expectedRepos := []Repository{}
	if !reflect.DeepEqual(repos, expectedRepos) {
		t.Fatalf("expected %v, but got %v", expectedRepos, repos)
	}
}

func TestGetPullRequests(t *testing.T) {
	s := NewGitHubService("dummyToken", "dummyOrg")

	prs, err := s.GetPullRequests(context.Background(), "testRepo")
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	expectedPRs := []PullRequest{}
	if !reflect.DeepEqual(prs, expectedPRs) {
		t.Fatalf("expected %v, but got %v", expectedPRs, prs)
	}
}

func TestGetReviews(t *testing.T) {
	s := NewGitHubService("dummyToken", "dummyOrg")

	reviews, err := s.GetReviews(context.Background(), "https://test.com/pr/1")
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	expectedReviews := []Review{}
	if !reflect.DeepEqual(reviews, expectedReviews) {
		t.Fatalf("expected %v, but got %v", expectedReviews, reviews)
	}
}
