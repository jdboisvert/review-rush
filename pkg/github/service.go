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
const maxPerPage = 100 // GitHub's maximum allowed value per page

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
	baseURL, err := s.getRepoUrl(s.org)
	if err != nil {
		return nil, err
	}

	var allRepos []Repository
	page := 1

	for {
		url := fmt.Sprintf("%s?page=%d&per_page=%d", baseURL, page, maxPerPage)

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

		// If we have an empty page, then we've fetched all repositories
		if len(repos) == 0 {
			break
		}

		allRepos = append(allRepos, repos...)
		page++
	}

	return allRepos, nil
}

func (s *service) GetPullRequests(ctx context.Context, repo string) ([]PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls?state=all&sort=updated&direction=desc", apiURL, s.org, repo)
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

	return prsFromToday(prs), nil
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

// Checks if the given time is from today.
func isFromToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// Used to filter out reviews that are not from today
// TODO: make this a attribute of the calling function so it can be passed in but for now it's fine.
func reviewsFromToday(reviews []Review) []Review {
	todayReviews := []Review{}
	for _, review := range reviews {
		if isFromToday(review.SubmittedAt) {
			todayReviews = append(todayReviews, review)
		}
	}
	return todayReviews
}

// Used to filter out PRs that are not from today it is assumed that the PRs are sorted by date in descending order
func prsFromToday(prs []PullRequest) []PullRequest {
	todayPRs := []PullRequest{}
	for _, pr := range prs {
		if isFromToday(pr.UpdatedAt) {
			todayPRs = append(todayPRs, pr)
		} else {
			// If the PR is not from today then we can assume that the rest of the PRs are not from today
			break
		}
	}
	return todayPRs
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
