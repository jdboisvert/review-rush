package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jdboisvert/review-rush/pkg/github"
	"github.com/jdboisvert/review-rush/pkg/slack"
	"github.com/jdboisvert/review-rush/pkg/utils"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	githubToken := os.Getenv("GITHUB_TOKEN")
	orgName := os.Getenv("ORG_NAME")
	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")

	fmt.Println("Starting Review Rush")
	fmt.Println("Github Org: ", githubToken)
	fmt.Println("Slack Channel: ", slackChannel)

	githubService := github.NewGitHubService(githubToken, orgName)
	slackService := slack.NewSlackService(slackToken)
	channel := slackChannel

	repos, err := githubService.GetRepos(ctx)
	if err != nil {
		log.Fatalf("Error getting repos: %v", err)
	}

	reviewCounts := make(map[string]int)

	for _, repo := range repos {
		fmt.Println("Getting PRs for repo: ", repo.Name)
		prs, err := githubService.GetPullRequests(ctx, repo.Name)
		if err != nil {
			log.Printf("Error getting PRs for repo %s: %v", repo.Name, err)
			continue
		}

		for _, pr := range prs {
			fmt.Println("Getting reviews for PR: ", pr.URL)
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

	fmt.Println("Review Counts")
	fmt.Println(reviewCounts)

	message := utils.FormatMessage(reviewCounts)

	if err := slackService.PostMessage(ctx, channel, message); err != nil {
		log.Fatalf("Error posting to Slack: %v", err)
	}
}
