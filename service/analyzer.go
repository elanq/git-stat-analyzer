package service

import (
	"context"

	"github.com/elanq/git-stat-analyzer/model"
)

type Analyzer interface {
	Analyze(ctx context.Context, filepath string) error
	GetUserStats(ctx context.Context, email string, repository string) ([]*model.AggregatedStat, error)
	GetAllStats(ctx context.Context, repository string) ([]*model.AggregatedStat, error)
}
