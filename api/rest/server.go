package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Viva-Victoria/bear-dtm/log"
	"github.com/Viva-Victoria/bear-dtm/service"
	"github.com/gorilla/mux"
)

type Server struct {
	logger log.Logger
	server http.Server
}

const Timeout = time.Second * 5

func NewServer(addr string, l log.Logger, s service.Service) *Server {
	router := mux.NewRouter()

	router.Handle("/api/1.0/transaction", CreateHandler(l, s)).Methods(http.MethodPost)
	router.Handle("/api/1.0/transaction/{id}/action", AddActionHandler(l, s)).Methods(http.MethodPost)
	router.Handle("/api/1.0/transaction/{id}/{state}", UpdateStateHandler(l, s)).Methods(http.MethodPatch)

	return &Server{
		logger: l,
		server: http.Server{
			Addr:              addr,
			Handler:           router,
			ReadTimeout:       Timeout,
			WriteTimeout:      Timeout,
			ReadHeaderTimeout: Timeout,
		},
	}
}

func (s *Server) Start() {
	s.logger.Info("server listening...")
	err := s.server.ListenAndServe()
	if err != nil {
		s.logger.Warn(fmt.Sprintf("server stopped: %v", err))
	}
}

func (s *Server) StartAsync() {
	go func() {
		s.Start()
	}()
	s.logger.Info("server started")
}

func (s *Server) Stop(ctx context.Context) {
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Warn(fmt.Sprintf("can't shutdown server: %v", err))
		return
	}
	s.logger.Info("server gracefully shutdown")
}

func (s *Server) StopForce() {
	err := s.server.Close()
	if err != nil {
		s.logger.Warn(fmt.Sprintf("can't force stop server: %v", err))
		return
	}
	s.logger.Info("server stopped")
}
