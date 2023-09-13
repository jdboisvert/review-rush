package github

import "time"

type GitHubUser struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

type Repository struct {
	Name string `json:"name"`
}

type PullRequest struct {
	URL       string    `json:"url"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Review struct {
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	SubmittedAt time.Time `json:"submitted_at"`
}
