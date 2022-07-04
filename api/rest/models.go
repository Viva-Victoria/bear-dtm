package rest

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/Viva-Victoria/bear-dtm/models"
	"github.com/google/uuid"
)

type HttpAction struct {
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Headers    map[string][]string `json:"headers"`
	BodyBase64 string              `json:"bodyBase64"`
}

func (h HttpAction) Validate() error {
	var builder strings.Builder

	if len(h.Method) == 0 {
		builder.WriteString("HTTP method should not be empty\n")
	}
	if len(h.URL) == 0 {
		builder.WriteString("HTTP URL should not be empty\n")
	}

	if builder.Len() == 0 {
		return nil
	}

	return errors.New(builder.String())
}

func (h *HttpAction) Map() (*models.HttpAction, error) {
	var err error
	var body []byte

	if len(h.BodyBase64) > 0 {
		body, err = base64.RawURLEncoding.DecodeString(h.BodyBase64)
		if err != nil {
			return nil, err
		}
	}

	return &models.HttpAction{
		Method:  h.Method,
		URL:     models.UrlFormat(h.URL),
		Headers: h.Headers,
		Body:    body,
	}, nil
}

type Action struct {
	Time         time.Time   `json:"time"`
	Name         string      `json:"name"`
	HttpRollback *HttpAction `json:"httpRollback,omitempty"`
}

func (a Action) Validate() error {
	var result strings.Builder

	if a.Time.IsZero() {
		result.WriteString("time should not be empty\n")
	}
	if len(a.Name) == 0 {
		result.WriteString("name should not be empty\n")
	}

	switch {
	case a.HttpRollback != nil:
		err := a.HttpRollback.Validate()
		if err != nil {
			result.WriteString(err.Error())
		}
	}

	if result.Len() == 0 {
		return nil
	}

	return errors.New(result.String())
}

func (a Action) Map() (models.Action, error) {
	var result strings.Builder
	var err error
	var http *models.HttpAction

	switch {
	case a.HttpRollback != nil:
		http, err = a.HttpRollback.Map()
		if err != nil {
			result.WriteString(err.Error())
		}
	}

	if result.Len() == 0 {
		return models.Action{
			Http: http,
		}, nil
	}

	return models.Action{}, errors.New(result.String())
}

type Transaction struct {
	State models.State           `json:"state"`
	Id    uuid.UUID              `json:"id"`
	Name  string                 `json:"name"`
	Tags  map[string]interface{} `json:"tags"`
}

func (t Transaction) Validate() error {
	var result strings.Builder

	if len(t.Name) == 0 {
		result.WriteString("name should not be empty\n")
	}

	if result.Len() == 0 {
		return nil
	}

	return errors.New(result.String())
}

func (t Transaction) Map() (models.Transaction, error) {
	return models.Transaction{
		State: t.State,
		Id:    t.Id,
		Name:  t.Name,
		Tags:  t.Tags,
	}, nil
}
