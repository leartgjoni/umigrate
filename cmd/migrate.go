package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"path/filepath"
)

func init() {
	migrateCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path")

	rootCmd.AddCommand(migrateCmd)
}

var (
	cfgFile    string
	migrateCmd = &cobra.Command{
		Use:              "migrate",
		Short:            "migrate file or all",
		Long:             `Migrate stuff`,
		PersistentPreRun: initConfig,
		Run: func(cmd *cobra.Command, args []string) {
			dbHost := viper.GetString("DB_HOST")
			dbPort := viper.GetString("DB_PORT")
			dbUser := viper.GetString("DB_USER")
			dbName := viper.GetString("DB_NAME")
			dbPassword := viper.GetString("DB_PASSWORD")

			if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" || dbPassword == "" {
				LogErr("DB_HOST, DB_PORT, DB_USER, DB_NAME, DB_PASSWORD are required configs\n")
			}

			dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
			db, err := sql.Open("postgres", dbUrl)
			if err != nil {
				LogErr("%s\n", err)
			}
			defer func() {
				if err := db.Close(); err != nil {
					LogErr("%s\n", err)
				}
			}()

			if len(args) >= 1 && args[0] != "" {
				runNamedMigration(args[0], db)
			} else {
				runAllMigrations(db)
			}
		},
	}
)

// initConfig reads config file and sets it with viper
func initConfig(*cobra.Command, []string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		LogErr("config file path is required\n")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		LogInfo("using config file: %s\n", viper.ConfigFileUsed())
	}
}

// runAllMigrations loops through all migrations in the current folder and runs all the migrations that have not been run before
func runAllMigrations(db *sql.DB) {
	LogInfo("RUNNING ALL MIGRATIONS\n")

	createMigrationsTable(db)

	filenamePattern := fmt.Sprintf("./*.sql")
	files, err := filepath.Glob(filenamePattern)
	if err != nil {
		LogErr("%s\n", err)
	}

	for _, filePath := range files {
		filename := filepath.Base(filePath)

		if checkIfMigrated(filename, db) {
			LogNotice("SKIPPING %s\n", filename)
			continue
		}

		sqlQuery, fileErr := ioutil.ReadFile(fmt.Sprintf("./%s", filename))
		if fileErr != nil {
			LogErr("%s\n", fileErr)
		}

		execQuery(db, string(sqlQuery))
		execQuery(db, "INSERT INTO _migrations (migration) VALUES($1)", filename)
		LogSuccess("%s MIGRATED\n", filename)
	}

}

// runNamedMigration runs migration file  if not run before
func runNamedMigration(migrationName string, db *sql.DB) {
	LogInfo("RUNNING %s MIGRATION\n", migrationName)

	createMigrationsTable(db)

	sqlQuery, fileErr := ioutil.ReadFile(fmt.Sprintf("./%s", migrationName))
	if fileErr != nil {
		LogErr("%s\n", fileErr)
	}

	if checkIfMigrated(migrationName, db) {
		LogNotice("SKIPPING %s\n", migrationName)
		return
	}

	execQuery(db, string(sqlQuery))
	execQuery(db, "INSERT INTO _migrations (migration) VALUES($1)", migrationName)
	LogSuccess("%s MIGRATED\n", migrationName)
}

// execQuery executes sql queries
func execQuery(db *sql.DB, sqlQuery string, args ...interface{}) {
	_, dbErr := db.Exec(sqlQuery, args...)
	if dbErr != nil {
		LogErr("%s\n", dbErr)
	}
}

// createMigrationsTable creates _migrations table in db
func createMigrationsTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS _migrations(
		   id serial PRIMARY KEY,
		   migration VARCHAR (255) UNIQUE NOT NULL,
		   migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		LogErr("%s\n", err)
	}
}

// checkIfMigrated checks if migration has been run before
func checkIfMigrated(filename string, db *sql.DB) bool {
	row := db.QueryRow("SELECT COUNT(*) FROM _migrations WHERE migration = $1", filename)

	var migrated bool
	if err := row.Scan(&migrated); err != nil {
		LogErr("%s\n", err)
	}

	return migrated
}
