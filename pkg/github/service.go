package github

import "context"

type GithubService interface {
	GetRepos(ctx context.Context) ([]Repository, error)
	GetPullRequests(ctx context.Context, repo string) ([]PullRequest, error)
	GetReviews(ctx context.Context, prURL string) ([]Review, error)
}

type service struct {
	token string
	org   string
}

func NewGitHubService(token, org string) GithubService {
	return &service{
		token: token,
		org:   org,
	}
}

func (s *service) GetRepos(ctx context.Context) ([]Repository, error) {
	return []Repository{}, nil
}

func (s *service) GetPullRequests(ctx context.Context, repo string) ([]PullRequest, error) {
	return []PullRequest{}, nil
}

func (s *service) GetReviews(ctx context.Context, prURL string) ([]Review, error) {
	return []Review{}, nil
}
