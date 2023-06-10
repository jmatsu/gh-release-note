package types

import "time"

type PullRequest struct {
	Assignees   []GitHubActor `json:"assignees"`
	Author      GitHubActor   `json:"author"`
	BaseRefName string        `json:"baseRefName"`
	HeadRefName string        `json:"headRefName"`
	CreatedAt   time.Time     `json:"createdAt"`
	MergeCommit struct {
		Sha string `json:"oid"`
	} `json:"mergeCommit"`
	MergedAt  time.Time   `json:"mergedAt"`
	MergedBy  GitHubActor `json:"mergedBy"`
	Number    int         `json:"number"`
	Title     string      `json:"title"`
	UpdatedAt time.Time   `json:"updatedAt"`
}
