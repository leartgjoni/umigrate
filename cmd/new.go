package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
		//fmt.Println("NEW CMD", viper.Get("DB_NAME"), args[0])

		migrationName := args[0]
		baseFilename := fmt.Sprintf("%s_%s", time.Now().Format("2006-01-02T15:04:05"), migrationName)
		content := []byte("/* migration file */")

		// write migration
		upFilename := fmt.Sprintf("%s.sql", baseFilename)
		if err := ioutil.WriteFile(upFilename, content, 0644); err != nil {
			log.Fatal(err)
		}
	},
}