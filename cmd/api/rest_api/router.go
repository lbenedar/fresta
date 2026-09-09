package restapi

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *AppServer) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.ShowText)
	router.HandlerFunc(http.MethodGet, "/v0/when", app.WhenNextSession)
	router.HandlerFunc(http.MethodGet, "/v0/healthcheck", app.healthcheckHandler)

	return router
}
