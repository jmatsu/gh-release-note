package service

import (
	"context"
	"fmt"
	"github.com/jmatsu/gh-release-note/prompts"
	"github.com/jmatsu/gh-release-note/types"
	"os"
	"testing"
)

var chatGPTService ChatGPTService

func init() {
	token := os.Getenv("GH_RELEASE_NOTE_OPENAI_API_KEY")

	chatGPTService = NewChatGPTService(context.TODO(), token)
}

func Test_buildPrompt(t *testing.T) {
	args := []struct {
		templateText string
		prs          []*types.PullRequest
	}{
		{
			templateText: prompts.SimpleTxt,
			prs: []*types.PullRequest{
				{
					Title:  "this is a pr1 title",
					Number: 1,
				},
				{
					Title:  "this is a pr2 title",
					Number: 2,
				},
			},
		},
		{
			templateText: "Please say hello",
			prs:          []*types.PullRequest{},
		},
	}

	for i, arg := range args {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			template, err := buildPrompt(arg.prs, arg.templateText)

			if err != nil {
				t.Errorf("failed due to %v", err)
			} else {
				t.Logf("Evaluated a template: %s", template)
			}
		})
	}
}

func TestChatGPTServiceImpl_GetSingleAnswer(t *testing.T) {
	if !chatGPTService.Active() {
		return
	}

	args := []struct {
		templateText string
		prs          []*types.PullRequest
	}{
		{
			templateText: prompts.SimpleTxt,
			prs: []*types.PullRequest{
				{
					Title:  "this is a pr1 title",
					Number: 1,
				},
				{
					Title:  "this is a pr2 title",
					Number: 2,
				},
			},
		},
		{
			templateText: "Please say hello",
			prs:          []*types.PullRequest{},
		},
	}

	for i, arg := range args {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			resp, err := chatGPTService.GetSingleAnswer(arg.prs, arg.templateText)

			if err != nil {
				t.Errorf("failed due to %v", err)
			} else {
				t.Logf("Got an answer: %s", resp)
			}
		})
	}
}
