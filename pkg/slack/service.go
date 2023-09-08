package slack

import "context"

type SlackService interface {
	PostMessage(ctx context.Context, message string) error
}

type service struct {
	token string
	url   string
}

func NewSlackService(token, url string) SlackService {
	return &service{token: token, url: url}
}

func (s *service) PostMessage(ctx context.Context, message string) error {
	return nil
}
