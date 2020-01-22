package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate file or all",
	Long:  `Migrate stuff`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("NEW CMD", viper.Get("DB_NAME"), args[0])

		dbHost := viper.GetString("DB_HOST")
		dbPort := viper.GetString("DB_PORT")
		dbUser := viper.GetString("DB_USER")
		dbName := viper.GetString("DB_NAME")
		dbPassword := viper.GetString("DB_PASSWORD")

		if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" || dbPassword == "" {
			log.Fatal("DB_HOST, DB_PORT, DB_USER, DB_NAME, DB_PASSWORD are required configs")
		}

		dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost,  dbPort, dbUser, dbName, dbPassword)
		db, err := sql.Open("postgres", dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if len(args) >= 1 && args[0] != "" {
			runNamedMigration(args[0], db)
		} else {
			runAllMigrations(db)
		}
	},
}

// runAllMigrations loops through all migrations in the current folder and runs all the migrations that have not been run before
func runAllMigrations(db *sql.DB) {
	fmt.Printf(InfoColor, "RUNNING ALL MIGRATIONS\n")

	createMigrationsTable(db)

	filenamePattern := fmt.Sprintf("./*.sql")
	files, err := filepath.Glob(filenamePattern)
	if err != nil {
		log.Fatal(err)
	}

	for _, filePath := range files {
		filename := filepath.Base(filePath)

		if checkIfMigrated(filename, db) {
			fmt.Printf(NoticeColor, fmt.Sprintf("SKIPPING %s\n", filename))
			continue
		}

		sqlQuery, fileErr := ioutil.ReadFile(fmt.Sprintf("./%s", filename))
		if fileErr != nil {
			log.Fatal(fileErr)
		}

		execQuery(db, string(sqlQuery))
		execQuery(db, "INSERT INTO _migrations (migration) VALUES($1)", filename)
		fmt.Printf(NoticeColor, fmt.Sprintf("%s MIGRATED\n", filename))
	}

}

// runNamedMigration runs migration file  if not run before
func runNamedMigration(migrationName string, db *sql.DB) {
	fmt.Printf(InfoColor, fmt.Sprintf("RUNNING %s MIGRATION\n", migrationName))

	createMigrationsTable(db)

	sqlQuery, fileErr := ioutil.ReadFile(fmt.Sprintf("./%s", migrationName))
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	if checkIfMigrated(migrationName, db) {
		fmt.Printf(NoticeColor, fmt.Sprintf("SKIPPING %s\n", migrationName))
		return
	}

	execQuery(db, string(sqlQuery))
	execQuery(db, "INSERT INTO _migrations (migration) VALUES($1)", migrationName)
	fmt.Printf(NoticeColor, fmt.Sprintf("%s MIGRATED\n", migrationName))
}


// execQuery executes sql queries
func execQuery(db *sql.DB, sqlQuery string, args ...interface{}) {
	_, dbErr := db.Exec(sqlQuery, args...)
	if dbErr != nil {
		log.Fatal(dbErr)
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
		log.Fatal(err)
	}
}

// checkIfMigrated checks if migration has been run before
func checkIfMigrated(filename string, db *sql.DB) bool {
	row := db.QueryRow("SELECT COUNT(*) FROM _migrations WHERE migration = $1", filename)

	var migrated bool
	row.Scan(&migrated)

	return migrated
}