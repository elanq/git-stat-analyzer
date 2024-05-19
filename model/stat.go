package model

import (
	"time"
)

type Set map[string]bool

type Stat struct {
	Email       string
	Commit      string
	Timestamp   time.Time
	FileChanges []string
	AddedLine   int
	RemovedLine int
}

type AggregatedStat struct {
	Timestamp        time.Time
	Repository       string
	ChangedFiles     Set
	TotalAddedLine   int
	TotalFileChanges int
	TotalRemovedLine int
}

type RepoStat struct {
	EarliestTimestamp time.Time
}
