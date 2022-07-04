package service

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Viva-Victoria/bear-dtm/models"
	"github.com/google/uuid"
)

type Worker interface {
	Do() error
}

type HttpWorker struct {
	transactionId uuid.UUID
	client        http.Client
	action        models.HttpAction
}

func NewHttpWorker(client http.Client, id uuid.UUID, action models.HttpAction) HttpWorker {
	return HttpWorker{
		client:        client,
		transactionId: id,
		action:        action,
	}
}

func (h HttpWorker) Do() error {
	request, err := http.NewRequest(h.action.Method, h.action.URL.Format(h.transactionId), bytes.NewBuffer(h.action.Body))
	if err != nil {
		return err
	}

	response, err := h.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode > 200 && response.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("response not OK: %s", response.Status)
}
