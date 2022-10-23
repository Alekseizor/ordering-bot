package repository

import (
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
)

const (
	txTimeout = 15 * time.Second
)

type Repository struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

// New создает новый репозиторий для доступа к данными
func New(db *sqlx.DB) (*Repository, error) {
	return &Repository{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}, nil
}
