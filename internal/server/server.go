package server

import (
	"context"
	"fmt"
	"net/http"
	"wbl0/WB_Task_L0/internal/config"
	getorder "wbl0/WB_Task_L0/internal/server/handlers/get_order"
	getorders "wbl0/WB_Task_L0/internal/server/handlers/get_orders"
	saveorder "wbl0/WB_Task_L0/internal/server/handlers/save_order"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	srv http.Server
	cfg *config.Config
	ctx *context.Context
}

type HFuncList struct {
	OrderGetter  getorder.OrderGetter
	OrdersGetter getorders.OrdersGetter
	OrderSaver   saveorder.OrderSaver
}

func New(cfg *config.Config, hf HFuncList, ctx *context.Context) *Server {
	return &Server{
		srv: http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
			Handler:      createRouter(hf),
			ReadTimeout:  cfg.Server.Timeout,
			WriteTimeout: cfg.Server.Timeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		cfg: cfg,
		ctx: ctx,
	}
}

func (s *Server) Start() {
	s.srv.ListenAndServe()
}

func (s *Server) Stop() {
	s.srv.Shutdown(*s.ctx)
}

func createRouter(hf HFuncList) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)

	router.Get("/", getorders.New(hf.OrdersGetter))
	router.Get("/{order_uid}", getorder.New(hf.OrderGetter))
	router.Post("/", saveorder.New(hf.OrderSaver))

	return router
}
