package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/spy-cats", app.listSpyCatHandler)
	router.HandlerFunc(http.MethodPost, "/v1/spy-cats", app.createSpyCatHandler)
	router.HandlerFunc(http.MethodGet, "/v1/spy-cats/:id", app.getSpyCatHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/spy-cats/:id", app.deleteSpyCatHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/spy-cats/:id", app.updateSpyCatHandler)

	return app.recoverPanic(
		app.logRequestResponse(
			router,
		),
	)
}
