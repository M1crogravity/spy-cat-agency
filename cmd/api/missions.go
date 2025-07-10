package main

import (
	"errors"
	"net/http"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/service"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
)

type inputTarget struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

// @Summary Create a new mission
// @Description Create a new mission with targets
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mission body CreateMissionRequestDoc true "Mission Details"
// @Success 201 {object} MissionResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 422 {object} ValidationErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /missions [post]
func (app *application) createMissionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Targets []inputTarget `json:"targets"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	targets := make([]*model.Target, 0, len(input.Targets))
	for _, inputTarget := range input.Targets {
		target := &model.Target{
			Name:    inputTarget.Name,
			Country: inputTarget.Country,
			State:   model.Created,
		}
		targets = append(targets, target)
	}

	mission := &model.Mission{
		State:   model.Created,
		Targets: targets,
	}

	err = app.missionsService.CreateMission(r.Context(), mission)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTooFewTargets) || errors.Is(err, service.ErrTooMuchTargets):
			errs := make(map[string]string)
			errs["targets"] = err.Error()
			app.failedValidationResponse(w, r, errs)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusCreated, envelope{"mission": mission})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary List all missions
// @Description Get a list of all missions
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MissionsResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /missions [get]
func (app *application) listMissionHandler(w http.ResponseWriter, r *http.Request) {
	missions, err := app.missionsService.GetAll(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"missions": missions})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get a mission by ID
// @Description Get details of a specific mission by ID
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Mission ID"
// @Success 200 {object} MissionResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /missions/{id} [get]
func (app *application) getMissionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	mission, err := app.missionsService.GetMissionByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorModelNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, mission)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Delete a mission
// @Description Delete a mission by ID
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Mission ID"
// @Success 200 {object} MessageResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /missions/{id} [delete]
func (app *application) deleteMissionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.missionsService.RemoveMission(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorModelNotFound):
			app.notFoundResponse(w, r)
		default:
			app.badRequestResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"message": "mission successfully deleted"})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Complete a mission
// @Description Mark a mission as completed
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Mission ID"
// @Success 200 {object} MissionResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /missions/{id}/complete [patch]
func (app *application) completeMissionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	mission, err := app.missionsService.CompleteMission(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorModelNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return

	}

	err = app.writeJson(w, http.StatusOK, envelope{"mission": mission})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Assign mission to spy cat
// @Description Assign a mission to a specific spy cat
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Mission ID"
// @Param spy-cat-id path int true "Spy Cat ID"
// @Success 200 {object} MissionResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /missions/{id}/spy-cat/{spy-cat-id} [patch]
func (app *application) assignMissionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	spyCatId, err := app.readIDParam(r, "spy-cat-id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	mission, err := app.missionsService.GetMissionByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrorModelNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	spyCat, err := app.spyCatsService.GetById(r.Context(), spyCatId)
	if err != nil {
		if errors.Is(err, storage.ErrorModelNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.missionsService.AssignMission(r.Context(), mission, spyCat)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAlreadyAssigned) || errors.Is(err, service.ErrSpyCatIsBusy):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"mission": mission, "spy-cat": spyCat})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Create mission target
// @Description Add a new target to a mission
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Mission ID"
// @Param target body CreateTargetRequestDoc true "Target Details"
// @Success 201 {object} TargetResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /missions/{id}/targets [post]
func (app *application) createMissionTargetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input inputTarget
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	target := &model.Target{
		Name:    input.Name,
		Country: input.Country,
		State:   model.Created,
	}
	mission, err := app.missionsService.GetMissionByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrorModelNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.missionsService.AddTarget(r.Context(), mission, target)
	if err != nil {
		switch err {
		case service.ErrOperationNotAllowedOnCompleted:
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"mission": mission})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Complete mission target
// @Description Mark a mission target as completed
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Mission ID"
// @Param target-id path int true "Target ID"
// @Success 200 {object} TargetResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /missions/{id}/targets/{target-id}/complete [patch]
func (app *application) completeMissionTargetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	targetId, err := app.readIDParam(r, "target-id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	mission, err := app.missionsService.GetMissionByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrorModelNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	spyCat := app.contextGetSpyCat(r)
	err = app.missionsService.CompleteTarget(r.Context(), mission, targetId, spyCat)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorModelNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, service.ErrAccessDenied):
			app.errorResponse(w, r, http.StatusForbidden, err.Error())
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"mission": mission})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Update mission target
// @Description Update the notes of a mission target
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Mission ID"
// @Param target-id path int true "Target ID"
// @Param notes body UpdateTargetNotesRequestDoc true "Target Notes"
// @Success 200 {object} TargetResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /missions/{id}/targets/{target-id} [patch]
func (app *application) updateMissionTargetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	targetId, err := app.readIDParam(r, "target-id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Notes string `json:"notes"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	mission, err := app.missionsService.GetMissionByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorModelNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	spyCat := app.contextGetSpyCat(r)
	err = app.missionsService.UpdateNotes(r.Context(), mission, targetId, input.Notes, spyCat)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAccessDenied):
			app.errorResponse(w, r, http.StatusForbidden, err.Error())
		case errors.Is(err, storage.ErrorModelNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, service.ErrOperationNotAllowedOnCompleted):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"mission": mission})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Delete mission target
// @Description Delete a mission target by ID
// @Tags missions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Mission ID"
// @Param target-id path int true "Target ID"
// @Success 200 {object} MessageResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /missions/{id}/targets/{target-id} [delete]
func (app *application) deleteMissionTargetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	targetId, err := app.readIDParam(r, "target-id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	mission, err := app.missionsService.GetMissionByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorModelNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.missionsService.RemoveTarget(r.Context(), mission, targetId)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMissionTargetMissmatch):
			app.notFoundResponse(w, r)
		case errors.Is(err, service.ErrOperationNotAllowedOnCompleted):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"mission": mission})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
