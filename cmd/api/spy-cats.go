package main

import (
	"errors"
	"net/http"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/validator"
)

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
