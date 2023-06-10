package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/jmatsu/gh-release-note/types"
	"github.com/sashabaranov/go-openai"
	"text/template"
)

type ChatGPTService interface {
	Active() bool
	GetSingleAnswer(prs []*types.PullRequest, templateText string) (string, error)
}

type chatGPTServiceImpl struct {
	ctx    context.Context
	client *openai.Client
}

func NewChatGPTService(ctx context.Context, token string) ChatGPTService {
	var client *openai.Client

	if token != "" {
		client = openai.NewClient(token)
	}

	return &chatGPTServiceImpl{
		ctx:    ctx,
		client: client,
	}
}

func (s *chatGPTServiceImpl) Active() bool {
	return s.client != nil
}

func (s *chatGPTServiceImpl) GetSingleAnswer(prs []*types.PullRequest, templateText string) (string, error) {
	if !s.Active() {
		return "", errors.New("this functionality is disabled")
	}

	prompt, err := buildPrompt(prs, templateText)

	if err != nil {
		return "", errors.Join(errors.New("can not build a prompt"), err)
	}

	// TODO Refine the output through multiple message chain

	messages := make([]openai.ChatCompletionMessage, 0)

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	resp, err := s.client.CreateChatCompletion(s.ctx, openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		Messages:  messages,
		MaxTokens: 500,
	})

	if err != nil {
		return "", errors.Join(errors.New("openapi returned an error"), err)
	}

	return resp.Choices[0].Message.Content, nil
}

func buildPrompt(prs []*types.PullRequest, templateText string) (string, error) {
	t, err := template.New("release-note").Parse(templateText)

	if err != nil {
		return "", errors.Join(errors.New("prompt cannot be parsed"), err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, struct {
		Prs []*types.PullRequest
	}{
		Prs: prs,
	})

	if err != nil {
		return "", errors.Join(errors.New("template cannot be evaluated properly"), err)
	}

	return buf.String(), nil
}
