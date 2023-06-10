package service

import (
	"context"
	"fmt"
	"github.com/cli/go-gh/v2/pkg/repository"
	"os"
	"testing"
	"time"
)

var ghService GitHubService
var repo repository.Repository
var sha string

func init() {
	name := "jmatsu/gh-release-note"
	sha = "6b2dc14ae9b38459155f9aef0d8762fab63d3d49"

	if v := os.Getenv("GH_RELEASE_NOTE_REPO"); v != "" {
		name = v
	}

	repo, _ = repository.Parse(name)

	if v := os.Getenv("GH_RELEASE_NOTE_COMMIT_SHA"); v != "" {
		sha = v
	}

	ghService = NewGitHubService(context.TODO(), repo)
}

func TestGitHubServiceImpl_ListMergedPullRequests(t *testing.T) {
	since, _ := time.Parse(time.DateOnly, "2023-06-09")
	until := time.Now()

	args := []struct {
		since *time.Time
		until *time.Time
	}{
		{since: nil, until: nil},
		{since: &since, until: nil},
		{since: nil, until: &until},
		{since: &since, until: &until},
	}

	for i, arg := range args {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			prs, err := ghService.ListMergedPullRequests(ListMergedPullRequestsOption{
				Base:          "main",
				MergedAtSince: arg.since,
				MergedAtUntil: arg.until,
			})

			if err != nil {
				t.Fatalf("fails due to %v", err)
				return
			}

			t.Logf("returns %d pull requests", len(prs))

			if len(prs) > 0 && prs[0].Number == 0 {
				t.Fatalf("pull request number is not assigned")
			} else {
				t.Logf("Get #%d: %s", prs[0].Number, prs[0].Title)
			}
		})
	}
}

func TestGitHubServiceImpl_ListTags(t *testing.T) {
	tags, err := ghService.ListTags(10)

	if err != nil {
		t.Fatalf("fails due to %v", err)
	}

	t.Logf("returns %d tags", len(tags))

	if len(tags) > 0 && tags[0].Name == "" {
		t.Fatalf("tag is not assigned")
	}
}
