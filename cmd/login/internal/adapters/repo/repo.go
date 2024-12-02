package repo

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sipki-tech/database/migrations"
	"medods/cmd/login/internal/app"
)

var _ app.Repo = &Repo{}

type (
	Config struct {
		Postgres   Connector
		MigrateDir string `yaml:"migrate_dir"`
		Driver     string `yaml:"driver"`
	}

	Connector struct {
		ConnectionDSN string `yaml:"connection_dsn"`
	}

	Repo struct {
		sql *sqlx.DB
	}
)

// DSN implements method for create migrations
func (c Connector) DSN() (string, error) {
	return c.ConnectionDSN, nil
}

// New creates Repo
func New(ctx context.Context, cfg Config) (*Repo, error) {

	migrates, err := migrations.Parse(cfg.MigrateDir)
	if err != nil {
		return nil, fmt.Errorf("migration.Parse: %w", err)
	}

	err = migrations.Run(ctx, cfg.Driver, &cfg.Postgres, migrations.Up, migrates)
	if err != nil {
		return nil, fmt.Errorf("migration.Run: %w", err)
	}

	dsn, err := cfg.Postgres.DSN()
	if err != nil {
		return nil, fmt.Errorf("connector.DSN: %w", err)
	}

	db, err := sqlx.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	return &Repo{
		db,
	}, nil
}

func (r *Repo) Close() error {
	return r.sql.Close()
}

// SaveToken saves hash by UserID
func (r *Repo) SaveToken(ctx context.Context, refreshToken app.TokensTable) error {
	const query = `
          INSERT INTO tokens_table (user_id, hash_refresh_token, ip, expires_at) VALUES ($1, $2, $3, $4)
          ON CONFLICT (user_id) DO UPDATE 
          SET hash_refresh_token = $2, ip = $3, expires_at = $4, updated_at = now()`
	_, err := r.sql.ExecContext(ctx, query, refreshToken.UserID, refreshToken.Hash, refreshToken.IP, refreshToken.ExpiresAt)
	if err != nil {
		return fmt.Errorf("r.sql.ExecContext: %w", err)
	}

	return nil
}

// TokenByUserID receives token information by UserID
func (r *Repo) TokenByUserID(ctx context.Context, userID uuid.UUID) (*app.TokensTable, error) {
	var dbToken RefreshToken
	const query = `SELECT * FROM tokens_table WHERE user_id = $1`
	err := r.sql.GetContext(ctx, &dbToken, query, userID)
	if err != nil {
		return nil, fmt.Errorf("r.sql.GetContext: %w", err)
	}

	return &app.TokensTable{
		ID:        dbToken.ID,
		UserID:    dbToken.UserID,
		Hash:      []byte(dbToken.RefreshToken),
		IP:        dbToken.IP,
		CreatedAt: dbToken.CreatedAt,
		UpdatedAt: dbToken.UpdatedAt,
		ExpiresAt: dbToken.ExpiresAt,
	}, nil

}
