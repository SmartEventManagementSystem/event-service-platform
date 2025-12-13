package integration

import (
	"context"
	"os"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start Postgres container
	pgReq := tc.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "local",
			"POSTGRES_PASSWORD": "local",
			"POSTGRES_DB":       "eventserviceplatformtpltest",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}
	pgC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{ContainerRequest: pgReq, Started: true})
	if err != nil {
		// If docker isn't available, fallback to running tests as-is
		os.Exit(m.Run())
	}

	// Start Redis container
	redisReq := tc.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(60 * time.Second),
	}
	redisC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{ContainerRequest: redisReq, Started: true})
	if err != nil {
		_ = pgC.Terminate(ctx)
		os.Exit(m.Run())
	}

	// Set environment variables so existing test setup picks up containers
	pgHost, _ := pgC.Host(ctx)
	pgPort, _ := pgC.MappedPort(ctx, "5432")
	os.Setenv("DB_URL", pgHost)
	os.Setenv("DB_PORT", pgPort.Port())
	os.Setenv("DB_NAME", "eventserviceplatformtpltest")
	os.Setenv("DB_USERNAME", "local")
	os.Setenv("DB_PASSWORD", "local")

	redisHost, _ := redisC.Host(ctx)
	redisPort, _ := redisC.MappedPort(ctx, "6379")
	os.Setenv("REDIS_HOSTS", redisHost+":"+redisPort.Port())

	// Run tests
	code := m.Run()

	// Teardown
	_ = redisC.Terminate(ctx)
	_ = pgC.Terminate(ctx)

	os.Exit(code)
}
