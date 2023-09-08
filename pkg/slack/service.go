package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const slackAPIEndpoint = "https://slack.com/api/chat.postMessage"

type SlackService interface {
	PostMessage(ctx context.Context, channel, message string) error
}

type service struct {
	client *http.Client
	token  string
	url    string
}

func NewSlackService(token, url string) SlackService {
	return &service{token: token, url: url, client: &http.Client{}}
}

func (s *service) PostMessage(ctx context.Context, channel, message string) error {
	payload := map[string]interface{}{
		"channel": channel,
		"text":    message,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", slackAPIEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to post message: %s", resp.Status)
	}

	var result struct {
		OK    bool   `json:"ok"`
		Error string `json:"error,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if !result.OK {
		return fmt.Errorf("failed to post message to Slack: %s", result.Error)
	}

	return nil
}
