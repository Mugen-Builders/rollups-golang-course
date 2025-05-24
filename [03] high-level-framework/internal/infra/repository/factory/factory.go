package factory

import (
	"fmt"
	"strings"

	. "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository/sqlite"
)

// NewRepositoryFromConnectionString chooses the backend based on the connection string.
// For instance:
//   - "postgres://user:pass@localhost/dbname" => Postgres
//   - "sqlite://some/path.db" => SQLite
//
// Then it initializes the repo, runs migrations, and returns it.
func NewRepositoryFromConnectionString(conn string) (Repository, error) {
	lowerConn := strings.ToLower(conn)
	switch {
	case strings.HasPrefix(lowerConn, "sqlite://"):
		return newSQLiteRepository(conn)
	default:
		return nil, fmt.Errorf("unrecognized connection string format: %s", conn)
	}
}

func newSQLiteRepository(conn string) (Repository, error) {
	sqliteRepo, err := sqlite.NewSQLiteRepository(conn)
	if err != nil {
		return nil, err
	}

	return sqliteRepo, nil
}
