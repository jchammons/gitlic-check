package swgithub

import (
	"fmt"

	"github.com/google/go-github/github"
)

var weekdays [7]string

func formatRepoStats(org string, repo string, commits []*github.WeeklyCommitActivity, mods []*github.WeeklyStats) string {
	totalCommits := 0
	linesAdded := 0
	linesRemoved := 0
	mostCommitsOneDay := 0
	mostCommitsOneWeek := 0
	var mostCommittable string
	highestDayCounts := map[string]int{"Sunday": 0, "Monday": 0, "Tuesday": 0, "Wednesday": 0, "Thursday": 0, "Friday": 0, "Saturday": 0}

	for _, commit := range commits {
		totalCommits += *commit.Total

		if *commit.Total > 0 {
			highestDay := 0
			highestDayInd := 0
			weekTotal := 0
			for i, dayTotal := range commit.Days {
				weekTotal += dayTotal

				if dayTotal > highestDay {
					highestDayInd = i
				}
				if dayTotal > mostCommitsOneDay {
					mostCommitsOneDay = dayTotal
				}
			}
			highestDayCounts[weekdays[highestDayInd]]++
			if weekTotal > mostCommitsOneWeek {
				mostCommitsOneWeek = weekTotal
			}
		}
	}

	if totalCommits > 0 {
		highestCount := 0
		for day, count := range highestDayCounts {
			if count > highestCount {
				mostCommittable = day
			}
		}
	} else {
		mostCommittable = "N/A"
	}

	mostLinesAddedOneWeek := 0
	mostLinesRemovedOneWeek := 0
	for _, mod := range mods {
		linesAdded += *mod.Additions
		linesRemoved += *mod.Deletions

		if *mod.Additions > mostLinesAddedOneWeek {
			mostLinesAddedOneWeek = *mod.Additions
		}

		if *mod.Deletions < mostLinesRemovedOneWeek {
			mostLinesRemovedOneWeek = *mod.Deletions
		}
	}

	return fmt.Sprint(org, ",", repo, ",", totalCommits, ",", mostCommitsOneDay, ",", mostCommitsOneWeek, ",", mostCommittable, ",", linesAdded, ",", linesRemoved, ",", mostLinesAddedOneWeek, ",", mostLinesRemovedOneWeek, "\n")
}
