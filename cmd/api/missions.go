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
