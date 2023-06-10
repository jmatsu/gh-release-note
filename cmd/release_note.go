package cmd

import (
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/jmatsu/gh-release-note/prompts"
	"github.com/jmatsu/gh-release-note/service"
	"github.com/jmatsu/gh-release-note/types"
	"github.com/urfave/cli/v2"
	"strings"
	"time"
)

const (
	SimplePromptStyle PromptStyle = "simple"
	CustomPromptStyle PromptStyle = "custom"
)

type PromptStyle string

type ReleaseNoteOption struct {
	Repo           repository.Repository
	Base           string
	Limit          int
	SinceTagName   *string
	SinceDate      *time.Time
	UntilTagName   *string
	UntilDate      *time.Time
	SkipGenerate   bool
	OpenAIAPIToken string
	Prompt         string
	PromptStyle    PromptStyle
}

func GenerateReleaseNote(c *cli.Context, f func(c *cli.Context) (ReleaseNoteOption, error)) error {
	opts, err := f(c)

	if err != nil {
		return err
	}

	github := service.NewGitHubService(c.Context, opts.Repo)
	var chatgpt service.ChatGPTService

	if !opts.SkipGenerate {
		chatgpt = service.NewChatGPTService(c.Context, opts.OpenAIAPIToken)
	}

	sinceDate, untilDate := opts.SinceDate, opts.UntilDate

	if sinceDate == nil {
		tags, err := github.ListTags(opts.Limit)

		if err != nil {
			return errors.Join(errors.New("cannot get tags preserved on GitHub"), err)
		}

		if opts.SinceTagName != nil {
			for _, t := range tags {
				if t.Name == *opts.SinceTagName {
					sinceDate = &t.Commit.CommittedDate
					break
				}
			}

			return errors.New(fmt.Sprintf("%s is not found in this repository", *opts.SinceTagName))
		}

		if sinceDate == nil {
			if t, err := chooseTagViaPrompt("Complete --since-tag", tags); err != nil {
				return errors.Join(errors.New("cannot choose the since tag"), err)
			} else {
				sinceDate = &t.Commit.CommittedDate
			}
		}

		// Don't suggest until-tag.
	}

	prs, err := github.ListMergedPullRequests(service.ListMergedPullRequestsOption{
		Base:          opts.Base,
		Limit:         opts.Limit,
		MergedAtSince: sinceDate,
		MergedAtUntil: untilDate,
	})

	if err != nil {
		return errors.Join(errors.New(fmt.Sprintf("cannot get pull requests between %s and %s", sinceDate.Format(time.DateTime), untilDate.Format(time.DateTime))), err)
	}

	if len(prs) == 0 {
		fmt.Println("No pull request is found.")
		return nil
	}

	if opts.SkipGenerate {
		for _, pr := range prs {
			fmt.Printf("#%d %s\n", pr.Number, pr.Title)
		}
		return nil
	}

	if opts.PromptStyle == SimplePromptStyle {

	}

	var aiPrompt string

	switch opts.PromptStyle {
	case SimplePromptStyle:
		aiPrompt = prompts.SimpleTxt
		break
	case CustomPromptStyle:
		aiPrompt = opts.Prompt
		break
	default:
		return errors.New(fmt.Sprintf("%s is an unknown prompt style", opts.PromptStyle))
	}

	releaseNote, err := chatgpt.GetSingleAnswer(prs, aiPrompt)

	if err != nil {
		return errors.Join(errors.New("could get pull requests but failed to generate a release note"), err)
	}

	fmt.Println(releaseNote)

	return nil
}

func chooseTagViaPrompt(prefixPrompt string, tags []*types.GitTag) (*types.GitTag, error) {
	suggests := make([]prompt.Suggest, len(tags))

	for i, tag := range tags {
		suggests[i] = prompt.Suggest{
			Text: fmt.Sprintf("%s (%s commited at %s)", tag.Name, tag.Commit.AbbreviatedOid, tag.Commit.CommittedDate.Format(time.DateTime)),
		}
	}

	chosen := prompt.Input(fmt.Sprintf("%s > ", prefixPrompt), func(document prompt.Document) []prompt.Suggest {
		return prompt.FilterContains(suggests, document.GetWordBeforeCursor(), true)
	}, prompt.OptionShowCompletionAtStart(), prompt.OptionCompletionOnDown(), prompt.OptionMaxSuggestion(8))

	for _, t := range tags {
		if strings.HasPrefix(chosen, t.Name) {
			return t, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("the chosen value (%s) seems to be invalid so no tag is found", chosen))
}
