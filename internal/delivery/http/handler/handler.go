package handler

import (
	"net/http"

	"github.com/Be1chenok/levelZero/internal/config"
	appService "github.com/Be1chenok/levelZero/internal/service"
	"github.com/gorilla/mux"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

type Handler struct {
	service *appService.Service
	conf    *config.Config
}

func New(conf *config.Config, service *appService.Service) *Handler {
	return &Handler{
		service: service,
		conf:    conf,
	}
}

func (h Handler) InitRoutes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/order/{uid:[a-zA-Z0-9]+}", h.FindOrderByUID).Methods("GET")

	return router
}
