package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiURL = "https://api.github.com"

type GithubService interface {
	GetRepos(ctx context.Context) ([]Repository, error)
	GetPullRequests(ctx context.Context, repo string) ([]PullRequest, error)
	GetReviews(ctx context.Context, prURL string) ([]Review, error)
}

type service struct {
	client *http.Client
	token  string
	org    string
}

func NewGitHubService(token, org string) GithubService {
	return &service{
		token:  token,
		org:    org,
		client: &http.Client{},
	}
}

func (s *service) GetRepos(ctx context.Context) ([]Repository, error) {
	url := fmt.Sprintf("%s/orgs/%s/repos", apiURL, s.org)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var repos []Repository
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func (s *service) GetPullRequests(ctx context.Context, repo string) ([]PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls", apiURL, s.org, repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var prs []PullRequest
	if err := json.Unmarshal(body, &prs); err != nil {
		return nil, err
	}

	return prs, nil
}

func (s *service) GetReviews(ctx context.Context, prURL string) ([]Review, error) {
	req, err := http.NewRequest(http.MethodGet, prURL+"/reviews", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reviews []Review
	if err := json.Unmarshal(body, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}
