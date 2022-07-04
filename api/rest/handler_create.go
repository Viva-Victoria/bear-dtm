package rest

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Viva-Victoria/bear-dtm/log"
	"github.com/Viva-Victoria/bear-dtm/service"
)

func CreateHandler(l log.Logger, s service.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(r.Body)
		if err != nil {
			logError(l, w, http.StatusBadRequest, err, "can't read request body")
			return
		}

		var restTransaction Transaction
		err = json.NewDecoder(r.Body).Decode(&restTransaction)
		if err != nil {
			logError(l, w, http.StatusBadRequest, err, "can't parse rest.Transaction")
			return
		}

		err = restTransaction.Validate()
		if err != nil {
			logError(l, w, http.StatusBadRequest, err, "rest.Transaction is invalid")
			return
		}

		serviceTransaction, err := restTransaction.Map()
		if err != nil {
			logError(l, w, http.StatusBadRequest, err, "can't map rest.Transaction to service.Transaction")
			return
		}

		serviceTransaction, err = s.Create(serviceTransaction)
		if err != nil {
			logError(l, w, http.StatusInternalServerError, err, "can't create transaction")
			return
		}

		buffer.Reset()
		err = json.NewEncoder(buffer).Encode(serviceTransaction)
		if err != nil {
			logError(l, w, http.StatusInternalServerError, err, "can't write rest.Transaction to response")
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(buffer.Bytes())
	})
}
