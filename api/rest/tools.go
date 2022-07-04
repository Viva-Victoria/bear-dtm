package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Viva-Victoria/bear-dtm/log"
	"github.com/Viva-Victoria/bear-dtm/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func logError(l log.Logger, w http.ResponseWriter, status int, err error, message string) {
	l.Error(err, message)
	w.WriteHeader(status)
	_, _ = w.Write([]byte(fmt.Sprintf("%s: %v", message, err)))
}

func getId(r *http.Request) (uuid.UUID, error) {
	raw, ok := mux.Vars(r)["id"]
	if !ok {
		return uuid.UUID{}, errors.New("no id in request")
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid id")
	}

	return id, nil
}

func getState(r *http.Request) (models.State, error) {
	state, ok := mux.Vars(r)["state"]
	if !ok {
		return 0, errors.New("no state in request")
	}

	switch strings.TrimSpace(strings.ToLower(state)) {
	case "confirm":
		return models.StateConfirmed, nil
	case "fail":
		return models.StateFailed, nil
	default:
		return 0, fmt.Errorf("state \"%s\" is not allowed", state)
	}
}
