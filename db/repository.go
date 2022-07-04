package db

import (
	"errors"

	"github.com/Viva-Victoria/bear-dtm/models"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("transaction not found")
)

type Repository interface {
	Get(id uuid.UUID) (models.Transaction, error)
	Insert(transaction models.Transaction) error
	Update(transaction models.Transaction) error
}
