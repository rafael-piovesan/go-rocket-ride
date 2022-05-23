package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	// loading bun's official Postgres driver.
	_ "github.com/uptrace/bun/driver/pgdriver"
)

// NewPostgresContainer creates a Postgres container and returns its DSN to be used
// in tests along with a termination callback to stop the container.
func NewPostgresContainer() (string, func(context.Context) error, error) {
	ctx := context.Background()

	templateURL := "postgres://postgres:postgres@localhost:%s/testdb?sslmode=disable"

	// Create the container
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:14.1",
			ExposedPorts: []string{
				"0:5432",
			},
			Env: map[string]string{
				"POSTGRES_DB":       "testdb",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_SSL_MODE": "disable",
			},
			Cmd: []string{
				"postgres", "-c", "fsync=off",
			},
			WaitingFor: wait.ForSQL(
				"5432/tcp",
				"pg",
				func(p nat.Port) string {
					return fmt.Sprintf(templateURL, p.Port())
				},
			).Timeout(time.Second * 15),
		},
	})
	if err != nil {
		return "", func(context.Context) error { return nil }, err
	}

	// Find ports assigned to the new container
	ports, err := c.Ports(ctx)
	if err != nil {
		return "", func(context.Context) error { return nil }, err
	}

	// Format driverURL
	driverURL := fmt.Sprintf(templateURL, ports["5432/tcp"][0].HostPort)

	return driverURL, c.Terminate, nil
}
