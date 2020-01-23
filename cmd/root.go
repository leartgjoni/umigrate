package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	infoColor    = "\033[1;34m%s\033[0m"
	successColor = "\033[1;32m%s\033[0m"
	noticeColor  = "\033[1;36m%s\033[0m"
	errorColor   = "\033[1;31m%s\033[0m"
)

var (
	rootCmd = &cobra.Command{
		Use:   "umigrate",
		Short: "Micro migrate",
		Long:  `Micro migrate is an extra lightweight tool for migrations in go`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func LogErr(format string, a ...interface{}) {
	fmt.Printf(errorColor, fmt.Sprintf(format, a...))
	os.Exit(1)
}

func LogInfo(format string, a ...interface{}) {
	fmt.Printf(infoColor, fmt.Sprintf(format, a...))
}

func LogNotice(format string, a ...interface{}) {
	fmt.Printf(noticeColor, fmt.Sprintf(format, a...))
}

func LogSuccess(format string, a ...interface{}) {
	fmt.Printf(successColor, fmt.Sprintf(format, a...))
}