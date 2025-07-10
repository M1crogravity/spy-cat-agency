package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/spy-cats", app.requireAgent(app.listSpyCatHandler))
	router.HandlerFunc(http.MethodPost, "/v1/spy-cats", app.requireAgent(app.createSpyCatHandler))
	router.HandlerFunc(http.MethodGet, "/v1/spy-cats/:id", app.requireAgent(app.getSpyCatHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/spy-cats/:id", app.requireAgent(app.deleteSpyCatHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/spy-cats/:id", app.requireAgent(app.updateSpyCatHandler))

	router.HandlerFunc(http.MethodPost, "/v1/missions", app.requireAgent(app.createMissionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/missions", app.requireAgent(app.listMissionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/missions/:id", app.requireAgent(app.getMissionHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/missions/:id", app.requireAgent(app.deleteMissionHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/missions/:id/complete", app.requireAgent(app.completeMissionHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/missions/:id/spy-cat/:spy-cat-id", app.requireAgent(app.assignMissionHandler))
	router.HandlerFunc(http.MethodPost, "/v1/missions/:id/targets", app.requireAgent(app.createMissionTargetHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/missions/:id/targets/:target-id/complete", app.requireSpyCat(app.completeMissionTargetHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/missions/:id/targets/:target-id", app.requireSpyCat(app.updateMissionTargetHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/missions/:id/targets/:target-id", app.requireAgent(app.deleteMissionTargetHandler))

	router.HandlerFunc(http.MethodPost, "/v1/agents", app.createAgentHandler) //let it be public for demo

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication/spy-cats", app.createSpyCatAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication/agents", app.createAgentAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/swagger/*filepath", httpSwagger.WrapHandler)

	return app.recoverPanic(
		app.logRequestResponse(
			app.authenticate(
				router,
			),
		),
	)
}
