//go:build unit
// +build unit

package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	if os.Getenv("START_MAIN") == "1" {
		main()
		return
	}

	t.Run("Missing env vars", func(t *testing.T) {
		stdout, stderr, err := startSubprocess(t)
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			assert.Empty(t, stdout)
			assert.Contains(t, stderr, "cannot load config")
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	})

	t.Run("Invalid Postgres URL", func(t *testing.T) {
		stdout, stderr, err := startSubprocess(t, "DB_SOURCE=rides", "SERVER_ADDRESS=foo", "STRIPE_KEY=bar")
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			assert.Empty(t, stdout)
			assert.Contains(t, stderr, "cannot open database")
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	})

	t.Run("Invalid Server Address", func(t *testing.T) {
		stdout, _, err := startSubprocess(
			t,
			"STRIPE_MOCK_INIT_CHECK=false", // skip initial stripe mock check
			"DB_SOURCE=postgresql://usr:pass@localhost:5432/db?sslmode=disable",
			"SERVER_ADDRESS=0.0.0.0:as",
			"STRIPE_KEY=bar",
		)
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			assert.Contains(t, stdout, "error running server")
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	})
}

// startSubprocess calls "go test" command specifying the test target name "TestMain" and setting
// the env var "START_MAIN=1". It will cause the test to be run again, but this time calling the
// "main()" func. This way, it's possible to retrieve and inspect the app exit code along with
// the stdout and stderr as well.
// See more at: https://stackoverflow.com/a/33404435
func startSubprocess(t *testing.T, envs ...string) (string, string, error) {
	var cout, cerr bytes.Buffer

	// call test suit again specifying "TestMain" as the target
	// nolint
	cmd := exec.Command(os.Args[0], "-test.run=TestMain")

	// set "START_MAIN" env var along with any additional value provided as parameter
	envs = append(envs, "START_MAIN=1")
	cmd.Env = append(os.Environ(), envs...)

	// capture subprocess' stdout and stderr
	cmd.Stdout = &cout
	cmd.Stderr = &cerr

	err := cmd.Start()
	if err == nil {
		go func() {
			// kill the subprocess after given timeout
			time.Sleep(time.Second * 5)
			if err := cmd.Process.Kill(); err != nil {
				t.Logf("error trying to kill process: %v", err)
			}
		}()

		// wait for subprocess completion
		err = cmd.Wait()
	}

	stdout := cout.String()
	stderr := cerr.String()
	return stdout, stderr, err
}
