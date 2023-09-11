package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	url, err := s.getRepoUrl(s.org)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.client.Do(req)
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
	url := fmt.Sprintf("%s/repos/%s/%s/pulls?state=all", apiURL, s.org, repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.client.Do(req)
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
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := s.client.Do(req)
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

	return reviewsFromToday(reviews), nil
}

func reviewsFromToday(reviews []Review) []Review {
	// Used to filter out reviews that are not from today
	// TODO: make this a attribute of the calling function so it can be passed in but for now it's fine.
	todayReviews := []Review{}
	now := time.Now()
	for _, review := range reviews {
		if review.SubmittedAt.Year() == now.Year() &&
			review.SubmittedAt.Month() == now.Month() &&
			review.SubmittedAt.Day() == now.Day() {
			todayReviews = append(todayReviews, review)
		}
	}
	return todayReviews
}

func (s *service) getRepoUrl(repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", s.org)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting user:", err)
		return "", err
	}
	defer resp.Body.Close()

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		fmt.Println("Decode error:", err)
		return "", err
	}

	if user.Type == "Organization" {
		return fmt.Sprintf("%s/orgs/%s/repos", apiURL, s.org), nil
	} else {
		// If not an organization then it's assumed to always be a user
		return fmt.Sprintf("%s/users/%s/repos", apiURL, s.org), nil
	}
}
