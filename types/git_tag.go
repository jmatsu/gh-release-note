package types

import (
	"time"
)

type GitTag struct {
	Name   string    `json:"name"`
	Commit GitCommit `json:"target"`
}

type GitCommit struct {
	Oid            string    `json:"oid"`
	AbbreviatedOid string    `json:"abbreviatedOid"`
	CommittedDate  time.Time `json:"committedDate"`
}
