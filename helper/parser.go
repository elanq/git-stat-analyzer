package helper

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "Mon Jan 2 15:04:05 2006 -0700"

var regexAuthor = regexp.MustCompile(`\<(.*?)\>`)

func ParseAuthor(line string) string {
	ok := strings.HasPrefix(line, "Author")
	if !ok {
		return ""
	}

	res := regexAuthor.FindStringSubmatch(line)
	if len(res) == 0 {
		return ""
	}

	return res[1]

}

func ParseDate(line string) time.Time {
	ok := strings.HasPrefix(line, "Date:")
	if !ok {
		return time.Time{}
	}

	line, ok = strings.CutPrefix(line, "Date:")
	if !ok {
		return time.Time{}
	}

	line = strings.TrimSpace(line)

	t, err := time.Parse(TimeFormat, line)
	if err != nil {
		return time.Time{}
	}
	return t

}

func ParseCommit(line string) string {
	ok := strings.HasPrefix(line, "commit")
	if !ok {
		return ""
	}
	line, ok = strings.CutPrefix(line, "commit")
	if !ok {
		return ""
	}

	line = strings.TrimSpace(line)
	return line

}

func ParseCommitMessage(msg string) ([]string, int, int) {
	commitLines := strings.Split(string(msg), "\n")
	commitSummary := commitLines[0]
	commitSummaries := strings.Split(commitSummary, ",")
	if len(commitSummaries) == 0 {
		log.Println("ERROR: commit summaries size:", len(commitSummaries))
		return []string{}, 0, 0
	}

	var addedLine, removedLine int
	fileChanges := []string{}
	var sAddedLine, sRemovedLine = "0", "0"
	if len(commitSummaries) == 2 {
		if strings.HasSuffix(commitSummaries[1], " insertions(+)") {
			sAddedLine = strings.TrimSuffix(commitSummaries[1], " insertions(+)")
			sAddedLine = strings.TrimSpace(sAddedLine)
		}
		if strings.HasSuffix(commitSummaries[1], " deletions(-)") {
			sRemovedLine = strings.TrimSuffix(commitSummaries[1], " deletions(-)")
			sRemovedLine = strings.TrimSpace(sRemovedLine)
		}
	}

	if len(commitSummaries) == 3 {
		sAddedLine = strings.TrimSuffix(commitSummaries[1], " insertions(+)")
		sRemovedLine = strings.TrimSuffix(commitSummaries[2], " deletions(-)")
		sAddedLine = strings.TrimSpace(sAddedLine)
		sRemovedLine = strings.TrimSpace(sRemovedLine)
	}

	addedLine, err := strconv.Atoi(sAddedLine)
	if err != nil {
		addedLine = 0
	}
	removedLine, err = strconv.Atoi(sRemovedLine)
	if err != nil {
		removedLine = 0
	}

	for i := 1; i < len(commitLines)-1; i++ {
		files := strings.Split(commitLines[i], "|")
		if len(files) == 0 {
			continue
		}
		f := strings.TrimSpace(files[0])
		fileChanges = append(fileChanges, f)
	}

	return fileChanges, addedLine, removedLine
}
