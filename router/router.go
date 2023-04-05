package router

import (
	"net/http"

	"github.com/brilliant-monkey/frigate-notify/config"
	"github.com/brilliant-monkey/frigate-notify/controller"
	"github.com/gorilla/mux"
)

func CreateRouter(config *config.AppConfig) *mux.Router {
	router := mux.NewRouter()

	router.
		PathPrefix("/notify/v1").
		Path("/subscribe").
		Methods(http.MethodGet, http.MethodOptions).
		Handler(http.HandlerFunc(controller.Subscribe)).
		Name("subscriptions")

	return router
}
