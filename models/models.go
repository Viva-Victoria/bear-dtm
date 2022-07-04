package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type UrlFormat string

func (u UrlFormat) Format(id uuid.UUID) string {
	return strings.Replace(string(u), "{id}", id.String(), 1)
}

type HttpAction struct {
	Method  string
	URL     UrlFormat
	Headers map[string][]string
	Body    []byte
}

type Action struct {
	Time time.Time
	Name string
	Http *HttpAction
}

type State int

const (
	StatePending State = iota + 1
	StateConfirmed
	StateRolledBack
	StateFailed
)

func (s State) String() string {
	switch s {
	case StatePending:
		return "Pending"
	case StateConfirmed:
		return "Confirmed"
	case StateRolledBack:
		return "RolledBack"
	case StateFailed:
		return "Failed"
	default:
		return ""
	}
}

type Transaction struct {
	State State
	Id    uuid.UUID
	Name  string
	Tags  map[string]interface{}
	Steps []Action
}
