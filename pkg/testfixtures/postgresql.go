package testfixtures

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/go-testfixtures/testfixtures/v3"
)

// Load loads database test fixtures specified by the "data" param, a list of files and/or directories (paths
// should be relative to the project's root dir).
//
// It takes as input a Postgres DSN, a list of files/directories leading to the *.yaml files and one more
// optional param, which accepts a map with data to be used while parsing files looking for template
// placeholders to be replaced.
func Load(dsn string, data []string, tplData map[string]interface{}) error {
	if len(data) == 0 {
		return errors.New("list of fixtures files/directories is empty")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	// find out the absolute path to this file
	// it'll be used to determine the project's root path
	_, callerPath, _, _ := runtime.Caller(0) // nolint:dogsled

	// look for migrations source starting from project's root dir
	rootPath := fmt.Sprintf(
		"%s/../..",
		filepath.ToSlash(filepath.Dir(callerPath)),
	)

	// assemble a list of fixtures paths to be loaded
	for i := range data {
		data[i] = fmt.Sprintf("%v/%v", rootPath, filepath.ToSlash(data[i]))
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Template(),
		testfixtures.TemplateData(tplData),
		// Paths must come after Template() and TemplateData()
		testfixtures.Paths(data...),
	)
	if err != nil {
		return err
	}

	// load fixtures into DB
	return fixtures.Load()
}
