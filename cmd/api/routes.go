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

	router.HandlerFunc(http.MethodPost, "/v1/missions", app.createMissionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/missions", app.listMissionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/missions/:id", app.getMissionHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/missions/:id", app.deleteMissionHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/missions/:id/complete", app.completeMissionHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/missions/:id/spy-cat/:spy-cat-id", app.assignMissionHandler)
	router.HandlerFunc(http.MethodPost, "/v1/missions/:id/targets", app.createMissionTargetHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/missions/:id/targets/:target-id/complete", app.completeMissionTargetHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/missions/:id/targets/:target-id", app.updateMissionTargetHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/missions/:id/targets/:target-id", app.deleteMissionTargetHandler)

	return app.recoverPanic(
		app.logRequestResponse(
			app.authenticate(
				router,
			),
		),
	)
}
