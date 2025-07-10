package main

import (
	"errors"
	"net/http"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/validator"
)

// @Summary Create a new agent
// @Description Create a new agent with name and password
// @Tags agents
// @Accept json
// @Produce json
// @Param agent body CreateAgentRequestDoc true "Agent Details"
// @Success 201 {object} AgentResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 422 {object} ValidationErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /agents [post]
func (app *application) createAgentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Password string
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	agent := &model.Agent{
		Name: input.Name,
	}

	err = agent.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if model.ValidateAgent(v, agent); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.agentsService.Create(r.Context(), agent)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorUniqueConstraintViolation):
			v.AddError("name", "an agent with this name is already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusCreated, envelope{"agent": agent})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
