package postgres

import (
	"database/sql"
	"embed"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrations embed.FS

func InitDB(logger *zap.SugaredLogger, dbUrl string) (*sqlx.DB, error) {
	logger.Infof("Connecting to DB with URL: %s", dbUrl)
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "connect to database")
	}

	if err := runMigrations(logger, db.DB, "migrations"); err != nil {
		return nil, errors.WithMessage(err, "run migrations")
	}

	return db, nil
}

func runMigrations(logger *zap.SugaredLogger, db *sql.DB, migrationsDir string) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return errors.WithMessage(err, "set dialect")
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return errors.WithMessage(err, "apply migrations")
	}

	logger.Info("Migrations applied successfully")
	return nil
}
