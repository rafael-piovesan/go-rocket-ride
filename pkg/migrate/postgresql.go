package migrate

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// Load driver to read in migrations from the file system.
	// See: https://github.com/golang-migrate/migrate/tree/master/source/file
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Up executes all migrations found at the given source path against the
// database specified by given DSN.
func Up(dsn string, source string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// find out the absolute path to this file
	// it'll be used to determine the project's root path
	_, callerPath, _, _ := runtime.Caller(0) // nolint:dogsled

	// look for migrations source starting from project's root dir
	sourceURL := fmt.Sprintf(
		"file://%s/../../%s",
		filepath.ToSlash(filepath.Dir(callerPath)),
		filepath.ToSlash(source),
	)

	migration, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	return migration.Up()
}
