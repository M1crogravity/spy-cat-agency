package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/validator"
)

// @Summary Create spy cat authentication token
// @Description Authenticate a spy cat and return a JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body AuthenticationRequestDoc true "Spy Cat Credentials"
// @Success 201 {object} TokenResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 422 {object} ValidationErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /tokens/authentication/spy-cats [post]
func (app *application) createSpyCatAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	model.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	spyCat, err := app.spyCatsService.GetByName(r.Context(), input.Name)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorModelNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := spyCat.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.tokensService.Create(r.Context(), spyCat.Id, model.SpyCatUserType, 24*time.Hour, model.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusCreated, envelope{"authentication_token": token})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// @Summary Create agent authentication token
// @Description Authenticate an agent and return a JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body AuthenticationRequestDoc true "Agent Credentials"
// @Success 201 {object} TokenResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 422 {object} ValidationErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /tokens/authentication/agents [post]
func (app *application) createAgentAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	model.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	agent, err := app.agentsService.GetByName(r.Context(), input.Name)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorModelNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := agent.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.tokensService.Create(r.Context(), agent.Id, model.AgentUserType, 24*time.Hour, model.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusCreated, envelope{"authentication_token": token})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
