package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type DiffInfo struct {
	Diff        string
	StagedFiles []string
	Insertions  int
	Deletions   int
}

func GetStagedDiff() (*DiffInfo, error) {
	_, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	if err != nil {
		return nil, fmt.Errorf("not a git repository")
	}


	diff, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged diff: %w", err)
	}

	if len(diff) == 0 {
		return nil, fmt.Errorf("no staged changes found — run git add first")
	}


	filesOut, err := exec.Command(
		"git", "diff", "--cached", "--name-only",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	files := []string{}
	for _, f := range strings.Split(string(filesOut), "\n") {
		f = strings.TrimSpace(f)
		if f != "" {
			files = append(files, f)
		}
	}


	statOut, err := exec.Command(
		"git", "diff", "--cached", "--shortstat",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get diff stat: %w", err)
	}

	insertions, deletions := parseStat(string(statOut))

	return &DiffInfo{
		Diff:        string(diff),
		StagedFiles: files,
		Insertions:  insertions,
		Deletions:   deletions,
	}, nil
}

func Commit(message string) error {
	out, err := exec.Command("git", "commit", "-m", message).CombinedOutput()
if err != nil {
	return fmt.Errorf("commit failed: %s", string(out))
}
	return nil
}

func parseStat(stat string) (insertions int, deletions int) {
	parts := strings.Split(stat, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "insertion") {
			fmt.Sscanf(part, "%d", &insertions)
		}
		if strings.Contains(part, "deletion") {
			fmt.Sscanf(part, "%d", &deletions)
		}
	}
	return
}

func TruncateDiff(diff string, maxChars int) string {
	if len(diff) <= maxChars {
		return diff
	}
	return diff[:maxChars] + "\n... (diff truncated)"
}

type FileGroup struct {
	Files      []string
	Diff       string
	Insertions int
	Deletions  int
}

func GroupStagedFiles(diff *DiffInfo) []FileGroup {
	groups := make(map[string][]string)

	for _, file := range diff.StagedFiles {
		category := categorizeFile(file)
		groups[category] = append(groups[category], file)
	}


	if len(groups) <= 1 {
		return []FileGroup{
			{
				Files:      diff.StagedFiles,
				Diff:       diff.Diff,
				Insertions: diff.Insertions,
				Deletions:  diff.Deletions,
			},
		}
	}

	result := []FileGroup{}
	for _, files := range groups {
		result = append(result, FileGroup{
			Files: files,
		})
	}
	return result
}

func categorizeFile(file string) string {
	lower := strings.ToLower(file)

	switch {
	case strings.Contains(lower, "auth") ||
		strings.Contains(lower, "login") ||
		strings.Contains(lower, "jwt") ||
		strings.Contains(lower, "session"):
		return "auth"

	case strings.Contains(lower, "test") ||
		strings.Contains(lower, "_test.go") ||
		strings.Contains(lower, "spec"):
		return "tests"

	case strings.Contains(lower, "readme") ||
		strings.Contains(lower, ".md") ||
		strings.Contains(lower, "doc"):
		return "docs"

	case strings.Contains(lower, "docker") ||
		strings.Contains(lower, "ci") ||
		strings.Contains(lower, "yml") ||
		strings.Contains(lower, "yaml") ||
		strings.Contains(lower, "makefile"):
		return "config"

	case strings.Contains(lower, "ui") ||
		strings.Contains(lower, "component") ||
		strings.Contains(lower, "style") ||
		strings.Contains(lower, "css") ||
		strings.Contains(lower, "html"):
		return "ui"

	case strings.Contains(lower, "db") ||
		strings.Contains(lower, "database") ||
		strings.Contains(lower, "migration") ||
		strings.Contains(lower, "schema"):
		return "database"

	case strings.Contains(lower, "api") ||
		strings.Contains(lower, "route") ||
		strings.Contains(lower, "handler") ||
		strings.Contains(lower, "controller"):
		return "api"

	default:
		return "general"
	}
}

func GetFileDiff(file string) string {
	out, err := exec.Command(
		"git", "diff", "--cached", "--", file,
	).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func StageFiles(files []string) error {
	args := append([]string{"add"}, files...)
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stage files: %s", string(out))
	}
	return nil
}