package main

import (
	"context"
	"log"
	"time"

	"github.com/jdboisvert/review-rush/pkg/github"
	"github.com/jdboisvert/review-rush/pkg/slack"
	"github.com/jdboisvert/review-rush/pkg/utils"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO these will be env vars
	githubService := github.NewGitHubService("GITHUB_TOKEN", "ORG_NAME")
	slackService := slack.NewSlackService("SLACK_TOKEN", "SLACK_URL")
	channel := "SLACK_CHANNEL"

	repos, err := githubService.GetRepos(ctx)
	if err != nil {
		log.Fatalf("Error getting repos: %v", err)
	}

	reviewCounts := make(map[string]int)

	for _, repo := range repos {
		prs, err := githubService.GetPullRequests(ctx, repo.Name)
		if err != nil {
			log.Printf("Error getting PRs for repo %s: %v", repo.Name, err)
			continue
		}

		for _, pr := range prs {
			reviews, err := githubService.GetReviews(ctx, pr.URL)
			if err != nil {
				log.Printf("Error getting reviews for PR %s: %v", pr.URL, err)
				continue
			}

			for _, review := range reviews {
				reviewCounts[review.User.Login]++
			}
		}
	}

	message := utils.FormatMessage(reviewCounts)
	if err := slackService.PostMessage(ctx, channel, message); err != nil {
		log.Fatalf("Error posting to Slack: %v", err)
	}
}
