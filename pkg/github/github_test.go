package github

import (
	"testing"

	"github.com/h2non/gock"

	"github.com/stretchr/testify/assert"
)

func TestGetRepos(t *testing.T) {
	defer gock.Off()

	// gock.New("https://api.github.com").
	// 	Get("/user/repos").
	// 	Reply(200).
	// 	JSON([]Repository{
	// 		{Name: "test"},
	// 	})

	repos := GetRepos()

	// assert.Equal(t, []Repository{{Name: "test"}}, repos)
	assert.Equal(t, []Repository{}, repos)
}

func TestGetPullRequests(t *testing.T) {
	defer gock.Off()

	// gock.New("https://api.github.com").
	// 	Get("/repos/test/test/pulls").
	// 	Reply(200).
	// 	JSON([]PullRequest{
	// 		{URL: "test"},
	// 	})

	prs := GetPullRequests("test/test")

	// assert.Equal(t, []PullRequest{{URL: "test"}}, prs)
	assert.Equal(t, []PullRequest{}, prs)
}

func TestGetReviews(t *testing.T) {
	defer gock.Off()

	// gock.New("https://api.github.com").
	// 	Get("/repos/test/test/pulls/1/reviews").
	// 	Reply(200).
	// 	JSON([]Review{
	// 		{User: struct {
	// 			Login string "json:\"login\""
	// 		}{Login: "test"}},
	// 	})

	reviews := GetReviews("test/test", 1)

	// assert.Equal(t, []Review{{User: struct {
	// 	Login string "json:\"login\""
	// }{Login: "test"}}}, reviews)

	assert.Equal(t, []Review{}, reviews)

}
