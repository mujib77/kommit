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