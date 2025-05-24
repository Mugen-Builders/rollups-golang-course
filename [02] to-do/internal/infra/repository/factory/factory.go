package factory

import (
	"context"
	"fmt"
	"strings"

	. "github.com/henriquemarlon/cartesi-golang-series/to-do/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/to-do/internal/infra/repository/in_memory"
	"github.com/henriquemarlon/cartesi-golang-series/to-do/internal/infra/repository/sqlite"
)

// NewRepositoryFromConnectionString chooses the backend based on the connection string.
// For instance:
//   - "postgres://user:pass@localhost/dbname" => Postgres
//   - "sqlite://some/path.db" => SQLite
//
// Then it initializes the repo, runs migrations, and returns it.
func NewRepositoryFromConnectionString(ctx context.Context, conn string) (Repository, error) {
	lowerConn := strings.ToLower(conn)
	switch {
	case strings.HasPrefix(lowerConn, "memory://"):
		return newInMemoryRepository()
	case strings.HasPrefix(lowerConn, "sqlite://"):
		return newSQLiteRepository(ctx, conn)
	default:
		return nil, fmt.Errorf("unrecognized connection string format: %s", conn)
	}
}

func newInMemoryRepository() (Repository, error) {
	inMemoryRepo, err := in_memory.NewInMemoryRepository()
	if err != nil {
		return nil, err
	}

	return inMemoryRepo, nil
}

func newSQLiteRepository(ctx context.Context, conn string) (Repository, error) {
	sqliteRepo, err := sqlite.NewSQLiteRepository(ctx, conn)
	if err != nil {
		return nil, err
	}

	return sqliteRepo, nil
}
