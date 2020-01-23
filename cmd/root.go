package cmd

import (
	"github.com/spf13/cobra"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	SuccessColor = "\033[1;32m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
)

var (
	rootCmd = &cobra.Command{
		Use:   "umigrate",
		Short: "Micro migrate",
		Long: `Micro migrate is an extra lightweight tool for migrations in go`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}