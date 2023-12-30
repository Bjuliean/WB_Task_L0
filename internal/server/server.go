package server

import (
	"fmt"
	"net/http"
	"wbl0/WB_Task_L0/internal/config"
	getorder "wbl0/WB_Task_L0/internal/server/handlers/get_order"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	srv         http.Server
	cfg         *config.Config
}

type HFuncList struct {
	OrderGetter getorder.OrderGetter
}

func New(cfg *config.Config, hf HFuncList) *Server {
	return &Server{
		srv: http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
			Handler:      createRouter(hf),
			ReadTimeout:  cfg.Server.Timeout,
			WriteTimeout: cfg.Server.Timeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		cfg:         cfg,
	}
}

func (s *Server) Start() {
	s.srv.ListenAndServe()
}

func createRouter(hf HFuncList) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)

	router.Get("/{order_uid}", getorder.New(hf.OrderGetter))

	return router
}
