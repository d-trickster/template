package db

import (
	"context"
	"log/slog"
	"template/internal/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool

	log *slog.Logger
}

func NewPostgres(connString string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		pool: pool,
		log:  slog.Default().With(logging.ComponentAttr("postgres")),
	}, nil
}

func (p *Postgres) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := p.pool.QueryRow(ctx,
		"SELECT id FROM users WHERE id=$1",
		id,
	).Scan(&user.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		p.log.Error("Query raw failed", logging.ErrAttr(err))
		return nil, ErrInternal
	}
	return &user, nil
}

func (p *Postgres) CreateUser(ctx context.Context, user *User) error {
	_, err := p.pool.Exec(ctx,
		"INSERT INTO users VALUES ($1)",
		user.ID,
	)
	if err != nil {
		p.log.Error("Failed to insert", logging.ErrAttr(err))
		return ErrInternal
	}
	return nil
}
