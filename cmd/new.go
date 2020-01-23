package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "new migration file",
	Long:  `creates new sql file`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a migration name argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		baseFilename := fmt.Sprintf("%s_%s", time.Now().Format("2006-01-02T15:04:05"), migrationName)
		content := []byte("/* migration file */")

		// write migration
		upFilename := fmt.Sprintf("%s.sql", baseFilename)
		if err := ioutil.WriteFile(upFilename, content, 0644); err != nil {
			fmt.Printf(ErrorColor, fmt.Sprintf("%s\n", err))
			os.Exit(1)
		}
		fmt.Printf(SuccessColor, fmt.Sprintf("MIGRATION FILE CREATED: %s\n", baseFilename))
	},
}