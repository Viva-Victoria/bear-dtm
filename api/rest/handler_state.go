package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Viva-Victoria/bear-dtm/log"
	"github.com/Viva-Victoria/bear-dtm/models"
	"github.com/Viva-Victoria/bear-dtm/service"
)

func UpdateStateHandler(l log.Logger, s service.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(r)
		if err != nil {
			logError(l, w, http.StatusBadRequest, err, "")
			return
		}

		state, err := getState(r)
		if err != nil {
			logError(l, w, http.StatusBadRequest, err, "")
		}

		var serviceTransaction models.Transaction
		switch state {
		case models.StateConfirmed:
			serviceTransaction, err = s.Confirm(id)
			break
		case models.StateFailed:
			serviceTransaction, err = s.Confirm(id)
			break
		default:
			logError(l, w, http.StatusBadRequest, nil, fmt.Sprintf("state %d not allowed", state))
		}
		if err != nil {
			logError(l, w, http.StatusInternalServerError, err, "can't create transaction")
			return
		}

		buffer := new(bytes.Buffer)
		err = json.NewEncoder(buffer).Encode(serviceTransaction)
		if err != nil {
			logError(l, w, http.StatusInternalServerError, err, "can't write rest.Transaction to response")
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(buffer.Bytes())
	})
}
