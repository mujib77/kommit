package commit

import (
	"context"
	"fmt"
	

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

	groups := git.GroupStagedFiles(diff)

	provider, err := ai.New(cfg)
	if err != nil {
		return err
	}

	if len(groups) > 1 {
		fmt.Printf("\n  ◆ detected %d logical groups\n\n", len(groups))
		return handleAtomicCommit(groups, provider, cfg)
	}

	fmt.Println("  generating commit messages...")

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
		return nil
	}

	return git.Commit(message)
}

func handleAtomicCommit(
	groups []git.FileGroup,
	provider ai.Provider,
	cfg config.Config,
) error {
	fmt.Println("  generating messages for each group...")

	for i, group := range groups {
		diff := group.Diff
		if diff == "" {
			for _, file := range group.Files {
				diff += git.GetFileDiff(file)
			}
		}

		truncated := git.TruncateDiff(diff, 4000)
		msgs, err := provider.GenerateMessages(
			context.Background(),
			truncated,
			cfg.Style,
		)
		if err != nil {
			return err
		}

		if len(msgs) == 0 {
			continue
		}

		fmt.Printf("  Group %d — %d files\n", i+1, len(group.Files))
		for _, f := range group.Files {
			fmt.Printf("    %s\n", f)
		}
		fmt.Println()

		message, quit, err := ui.RunUI(
			msgs,
			group.Files,
			group.Insertions,
			group.Deletions,
		)
		if err != nil {
			return err
		}

		if quit {
			fmt.Println("  skipped")
			continue
		}

		err = git.StageFiles(group.Files)
		if err != nil {
			return err
		}

		err = git.Commit(message)
		if err != nil {
			return err
		}

		fmt.Printf("  ✓ committed group %d: %s\n\n", i+1, message)
	}

	return nil
}