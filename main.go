package main

import (
	"context"

	"github.com/elanq/git-stat-analyzer/service"
)

func main() {
	analyzer := service.NewGitAnalyzer()
	err := analyzer.Analyze(context.Background(), "/Users/elan/work/stock-ratio-engine")
	if err != nil {
		panic(err)
	}

	/*
		f, err := os.OpenFile("elan.qisthi@ajaib.co.id.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "    ")

		stats, err := analyzer.GetUserStats(context.Background(), "elan.qisthi@ajaib.co.id", "/Users/elan/work/stock-ratio-engine")
		if err != nil {
			panic(err)
		}

		if err := enc.Encode(stats); err != nil {
			panic(err)
		}
	*/
}
