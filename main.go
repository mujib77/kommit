package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/mujib77/kommit/cmd/commit"
	"github.com/mujib77/kommit/config"
)

func main() {
	root := &cobra.Command{
		Use:   "kommit",
		Short: "AI-powered git commit messages",
		Long:  `kommit generates meaningful git commit messages using AI.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return commit.NewCommitCmd().RunE(cmd, args)
		},
	}

	root.AddCommand(commit.NewCommitCmd())

	root.AddCommand(&cobra.Command{
		Use:   "setup",
		Short: "Show setup instructions",
		Run: func(cmd *cobra.Command, args []string) {
			config.PrintSetup()
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("kommit v0.1.0")
		},
	})

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}