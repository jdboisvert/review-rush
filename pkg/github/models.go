package github

type Repository struct {
	Name string `json:"name"`
}

type PullRequest struct {
	URL string `json:"url"`
}

type Review struct {
	User struct {
		Login string `json:"login"`
	} `json:"user"`
}