package service

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/Viva-Victoria/bear-dtm/db"
	"github.com/Viva-Victoria/bear-dtm/models"
	"github.com/google/uuid"
)

type Service interface {
	Create(transaction models.Transaction) (models.Transaction, error)
	AddAction(transactionId uuid.UUID, action models.Action) (models.Transaction, error)
	Confirm(transactionId uuid.UUID) (models.Transaction, error)
	Rollback(transactionId uuid.UUID) (models.Transaction, error)
}

type ServiceImpl struct {
	httpClient http.Client
	repository db.Repository
}

func (s ServiceImpl) Create(transaction models.Transaction) (models.Transaction, error) {
	transaction.Id = uuid.New()
	err := s.repository.Insert(transaction)
	if err != nil {
		return models.Transaction{}, err
	}

	return s.repository.Get(transaction.Id)
}

func (s ServiceImpl) AddAction(transactionId uuid.UUID, action models.Action) (models.Transaction, error) {
	transaction, err := s.repository.Get(transactionId)
	if err != nil {
		return models.Transaction{}, err
	}

	transaction.Steps = append(transaction.Steps, action)
	err = s.repository.Update(transaction)
	if err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

func (s ServiceImpl) Confirm(transactionId uuid.UUID) (models.Transaction, error) {
	transaction, err := s.repository.Get(transactionId)
	if err != nil {
		return models.Transaction{}, err
	}

	if transaction.State != models.StatePending {
		return models.Transaction{}, fmt.Errorf("can't confirm transaction at %s state", transaction.State)
	}

	transaction.State = models.StateConfirmed
	return transaction, s.repository.Update(transaction)
}

func (s ServiceImpl) Rollback(transactionId uuid.UUID) (models.Transaction, error) {
	transaction, err := s.repository.Get(transactionId)
	if err != nil {
		return models.Transaction{}, err
	}

	if transaction.State != models.StatePending {
		return models.Transaction{}, fmt.Errorf("can't rollback transaction at %s state", transaction.State)
	}

	steps := transaction.Steps
	sort.SliceStable(steps, func(i, j int) bool {
		return steps[i].Time.Before(steps[j].Time)
	})

	workers := make([]Worker, len(steps))
	for i, step := range steps {
		switch {
		case step.Http != nil:
			workers[i] = NewHttpWorker(s.httpClient, transaction.Id, *step.Http)
		}
	}

	for _, worker := range workers {
		err := worker.Do()
		if err != nil {
			// TODO: SHIT WE CAN'T FULLY ROLLBACK TRANSACTION, NOOO!!!
		}
	}

	transaction.State = models.StateRolledBack
	return transaction, s.repository.Update(transaction)
}
