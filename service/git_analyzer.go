package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/elanq/git-stat-analyzer/helper"
	"github.com/elanq/git-stat-analyzer/model"
)

type database map[string]*model.AggregatedStat
type repoStat map[string]*model.RepoStat
type rawStat map[string][]model.Stat
type users map[string]bool

type gitAnalyzer struct {
	db        database
	repoStats repoStat
	rawStat   rawStat
	users     users
}

const cmdGitHistory = "cd %s && git log > temp_log"

func timestampUserRepoIndex(timestamp time.Time, email, repository string) string {
	tf := timestamp.Format(time.DateOnly)
	return strings.Join([]string{tf, email, repository}, ":")
}

func (g *gitAnalyzer) Analyze(ctx context.Context, repositoryPath string) error {
	if _, err := os.Stat(repositoryPath); err == os.ErrNotExist {
		return err
	}
	defer removeTemp(repositoryPath)

	commands := []string{"cd", repositoryPath, "&&", "git", "log", ">", "temp_log"}
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", strings.Join(commands, " "))
	if err := cmd.Run(); err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/temp_log", repositoryPath)
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	stats := make([]model.Stat, 0)
	stat := model.Stat{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if stat.Commit != "" && stat.Email != "" && !stat.Timestamp.IsZero() {
			stats = append(stats, stat)
			if _, ok := g.users[stat.Email]; !ok {
				g.users[stat.Email] = true
			}

			if _, ok := g.rawStat[stat.Email]; ok {
				g.rawStat[stat.Email] = append(g.rawStat[stat.Email], stat)
			} else {
				g.rawStat[stat.Email] = make([]model.Stat, 0)
				g.rawStat[stat.Email] = append(g.rawStat[stat.Email], stat)
			}

			// aggregated stat db
			key := timestampUserRepoIndex(stat.Timestamp, stat.Email, repositoryPath)
			if _, ok := g.db[key]; ok {
				g.db[key].TotalAddedLine += stat.AddedLine
				g.db[key].TotalRemovedLine += stat.RemovedLine
				for _, f := range stat.FileChanges {
					if _, ok := g.db[key].ChangedFiles[f]; !ok {
						g.db[key].ChangedFiles[f] = true
					}
				}
				g.db[key].TotalFileChanges = len(g.db[key].ChangedFiles)
			} else {
				set := make(model.Set, 0)
				for _, f := range stat.FileChanges {
					if _, ok := set[f]; !ok {
						set[f] = true
					}
				}
				g.db[key] = &model.AggregatedStat{
					Repository:       repositoryPath,
					ChangedFiles:     set,
					TotalAddedLine:   stat.AddedLine,
					TotalFileChanges: len(set),
					TotalRemovedLine: stat.RemovedLine,
				}
			}
			stat = model.Stat{}
		}
		line := scanner.Text()
		if email := helper.ParseAuthor(line); email != "" {
			stat.Email = email
			continue
		}
		if commit := helper.ParseCommit(line); commit != "" {
			log.Println("analyze commit ", commit)
			stat.Commit = commit
			stat.FileChanges, stat.AddedLine, stat.RemovedLine = analyzeCommit(ctx, repositoryPath, commit)
			continue
		}
		if timestamp := helper.ParseDate(line); !timestamp.IsZero() {
			stat.Timestamp = timestamp
			continue
		}
	}
	lastStat := stats[len(stats)-1]
	if _, ok := g.repoStats[repositoryPath]; !ok {
		g.repoStats[repositoryPath] = &model.RepoStat{
			EarliestTimestamp: lastStat.Timestamp,
		}
	}

	return nil
}

func analyzeCommit(ctx context.Context, repositoryPath string, commit string) ([]string, int, int) {
	commands := []string{"cd", repositoryPath, "&&", "git", "--no-pager", "show", commit, "--stat", `--format=''`, "|", "tail", "-r"}
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", strings.Join(commands, " "))
	out, err := cmd.Output()
	if err != nil {
		return []string{}, 0, 0
	}
	return helper.ParseCommitMessage(string(out))
}

func removeTemp(repositoryPath string) error {
	filePath := fmt.Sprintf("%s/temp_log", repositoryPath)
	return os.Remove(filePath)
}

func (g *gitAnalyzer) GetUserStats(ctx context.Context, email string, repository string) ([]*model.AggregatedStat, error) {
	if _, ok := g.repoStats[repository]; !ok {
		return nil, errors.New(fmt.Sprintf("repository %v not found", repository))
	}

	results := make([]*model.AggregatedStat, 0)
	cursor := g.repoStats[repository].EarliestTimestamp
	for cursor.Before(time.Now()) {
		key := timestampUserRepoIndex(cursor, email, repository)
		if _, ok := g.db[key]; !ok {
			cursor = cursor.Add(24 * time.Hour)
			continue
		}

		res := *g.db[key]
		res.Timestamp = cursor
		results = append(results, &res)
		cursor = cursor.Add(24 * time.Hour)
	}
	return results, nil
}

func (g *gitAnalyzer) GetAllStats(ctx context.Context, repository string) ([]*model.AggregatedStat, error) {
	return nil, nil

}

func NewGitAnalyzer() Analyzer {
	return &gitAnalyzer{
		db:        make(database),
		repoStats: make(repoStat),
		rawStat:   make(rawStat),
		users:     make(users),
	}
}
