package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/jmatsu/gh-release-note/queries"
	"github.com/jmatsu/gh-release-note/types"
	"math"
	"strconv"
	"time"
)

const (
	// https://docs.github.com/en/search-github/getting-started-with-searching-on-github/understanding-the-search-syntax#query-for-dates
	queryDateFormat = "2006-01-02T15:04:05Z"
)

type GitHubService interface {
	ListTags(limit int) ([]*types.GitTag, error)
	ListMergedPullRequests(option ListMergedPullRequestsOption) ([]*types.PullRequest, error)
}

type gitHubServiceImpl struct {
	ctx  context.Context
	repo repository.Repository
}

func NewGitHubService(ctx context.Context, repo repository.Repository) GitHubService {
	return &gitHubServiceImpl{
		ctx:  ctx,
		repo: repo,
	}
}

func (c *gitHubServiceImpl) ListTags(limit int) ([]*types.GitTag, error) {
	opts := []string{"api", "graphql"}

	// repo control
	opts = append(opts, "-F", fmt.Sprintf("owner=%s", c.repo.Owner))
	opts = append(opts, "-F", fmt.Sprintf("name=%s", c.repo.Name))
	opts = append(opts, "-F", fmt.Sprintf("limit=%d", limit))

	opts = append(opts, "-f", fmt.Sprintf("query=%s", queries.ListTagGraphQL))

	opts = append(opts, "--jq", ".data.repository.refs.nodes")

	var tags []*types.GitTag

	stdOut, _, err := gh.ExecContext(c.ctx, opts...)

	if err != nil {
		return nil, errors.Join(errors.New("failed to list git tags"), err)
	} else if err := json.Unmarshal(stdOut.Bytes(), &tags); err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal git tags"), err)
	}

	return tags, nil
}

type ListMergedPullRequestsOption struct {
	Base          string
	Limit         int
	MergedAtSince *time.Time
	MergedAtUntil *time.Time
}

func (c *gitHubServiceImpl) ListMergedPullRequests(option ListMergedPullRequestsOption) ([]*types.PullRequest, error) {
	var prs []*types.PullRequest

	opts := []string{"pr", "list"}

	// repo control
	opts = append(opts, "--repo", fmt.Sprintf("%s/%s", c.repo.Owner, c.repo.Name))

	// size control
	limit := int64(math.Min(100, math.Max(float64(option.Limit), 5)))
	opts = append(opts, "--limit", strconv.FormatInt(limit, 10))

	// search criteria
	opts = append(opts, "--base", option.Base)

	since, until := "*", "*"

	if option.MergedAtSince != nil && option.MergedAtUntil != nil {
		since = option.MergedAtSince.UTC().Format(queryDateFormat)
		until = option.MergedAtUntil.UTC().Format(queryDateFormat)
	} else if option.MergedAtSince != nil {
		since = option.MergedAtSince.UTC().Format(queryDateFormat)
	} else if option.MergedAtUntil != nil {
		until = option.MergedAtUntil.UTC().Format(queryDateFormat)
	}

	if since != "*" || until != "*" {
		// https://docs.github.com/en/search-github/getting-started-with-searching-on-github/understanding-the-search-syntax#query-for-dates
		opts = append(opts, "--search", fmt.Sprintf("merged:%s..%s", since, until))
	} else {
		opts = append(opts, "is:merged")
	}

	// serialized attributes
	opts = append(opts, "--json", "number,title,author,baseRefName,mergedAt,updatedAt,milestone,mergedBy,createdAt,mergeCommit,assignees")

	stdOut, _, err := gh.ExecContext(c.ctx, opts...)

	if err != nil {
		return nil, errors.Join(errors.New("failed to list pull requests"), err)
	} else if err := json.Unmarshal(stdOut.Bytes(), &prs); err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal pull requests"), err)
	}

	return prs, nil
}
