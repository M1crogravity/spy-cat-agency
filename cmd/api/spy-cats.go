package main

import (
	"errors"
	"net/http"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/validator"
)

// @Summary List all spy cats
// @Description Get a list of all spy cats
// @Tags spy-cats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SpyCatsResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /spy-cats [get]
func (app *application) listSpyCatHandler(w http.ResponseWriter, r *http.Request) {
	//no paging
	spyCats, err := app.spyCatsService.GetAll(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"spy-cats": spyCats})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get a spy cat by ID
// @Description Get details of a specific spy cat by ID
// @Tags spy-cats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Spy Cat ID"
// @Success 200 {object} SpyCatResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /spy-cats/{id} [get]
func (app *application) getSpyCatHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	spyCat, err := app.spyCatsService.GetById(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorModelNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"spy-cat": spyCat})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Create a new spy cat
// @Description Create a new spy cat with the provided details
// @Tags spy-cats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param spy-cat body CreateSpyCatRequestDoc true "Spy Cat Details"
// @Success 201 {object} SpyCatResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 422 {object} ValidationErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /spy-cats [post]
func (app *application) createSpyCatHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name              string  `json:"name"`
		YearsOfExperience int     `json:"years_of_experience"`
		Breed             string  `json:"breed"`
		Salary            float64 `json:"salary"`
		Password          string  `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	spyCat := &model.SpyCat{
		Name:              input.Name,
		YearsOfExperience: input.YearsOfExperience,
		Breed:             input.Breed,
		Salary:            input.Salary,
	}

	err = spyCat.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	breeds, err := app.spyCatsService.GetBreeds(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if model.ValidateSpyCat(v, spyCat, breeds); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.spyCatsService.Create(r.Context(), spyCat)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorUniqueConstraintViolation):
			v.AddError("name", "a spy cat with this name is already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusCreated, envelope{"spy-cat": spyCat})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Delete a spy cat
// @Description Delete a spy cat by ID
// @Tags spy-cats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Spy Cat ID"
// @Success 200 {object} MessageResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /spy-cats/{id} [delete]
func (app *application) deleteSpyCatHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.spyCatsService.Remove(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrorModelNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJson(w, http.StatusOK, envelope{"message": "spy cat successfully deleted"})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Update spy cat salary
// @Description Update the salary of a specific spy cat
// @Tags spy-cats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Spy Cat ID"
// @Param salary body UpdateSpyCatSalaryRequestDoc true "Salary Update"
// @Success 200 {object} SpyCatResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Failure 404 {object} ErrorResponseDoc
// @Failure 500 {object} ErrorResponseDoc
// @Router /spy-cats/{id} [patch]
func (app *application) updateSpyCatHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Salary float64 `json:"salary"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	spyCat, err := app.spyCatsService.UpdateSalary(r.Context(), id, input.Salary)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"spy-cat": spyCat})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
