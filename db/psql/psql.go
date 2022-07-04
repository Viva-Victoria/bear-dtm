package psql

import (
	_ "embed"

	"github.com/Viva-Victoria/bear-dtm/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(connectionString string) (*Repository, error) {
	connConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	nativeDB := stdlib.OpenDB(*connConfig)
	db := sqlx.NewDb(nativeDB, "pgx")

	return &Repository{
		db: db,
	}, nil
}

func (r Repository) Get(id uuid.UUID) (models.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) Insert(transaction models.Transaction) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) Update(transaction models.Transaction) error {
	//TODO implement me
	panic("implement me")
}
