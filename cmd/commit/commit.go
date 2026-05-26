package commit

import (
	"context"
	"fmt"
	"os"

	"github.com/mujib77/kommit/config"
	"github.com/mujib77/kommit/internal/ai"
	"github.com/mujib77/kommit/internal/git"
	"github.com/mujib77/kommit/internal/ui"
	"github.com/spf13/cobra"
)

func NewCommitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Generate AI commit messages",
		RunE:  runCommit,
	}
	return cmd
}

func runCommit(cmd *cobra.Command, args []string) error {
	cfg := config.Load()

	if cfg.APIKey == "" {
		config.PrintSetup()
		return fmt.Errorf("no API key configured")
	}

	fmt.Println("\n  ◆ KOMMIT — analyzing your changes...")

	diff, err := git.GetStagedDiff()
	if err != nil {
		return err
	}

	fmt.Printf("  staged: %d files  +%d  -%d\n",
		len(diff.StagedFiles),
		diff.Insertions,
		diff.Deletions,
	)
	fmt.Println("  generating commit messages...")

	provider, err := ai.New(cfg)
	if err != nil {
		return err
	}

	truncated := git.TruncateDiff(diff.Diff, 8000)

	messages, err := provider.GenerateMessages(
		context.Background(),
		truncated,
		cfg.Style,
	)
	if err != nil {
		return err
	}

	message, quit, err := ui.RunUI(
		messages,
		diff.StagedFiles,
		diff.Insertions,
		diff.Deletions,
	)
	if err != nil {
		return err
	}

	if quit {
		fmt.Println("  cancelled")
		return nil
	}

	err = git.Commit(message)
	if err != nil {
		return err
	}

	fmt.Printf("\n  ✓ committed: %s\n\n", message)
	os.Exit(0)
	return nil
}