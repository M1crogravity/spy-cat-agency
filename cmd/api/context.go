package main

import (
	"context"
	"net/http"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
)

type contextKey string

const spyCatContextKey = contextKey("spy-cat")

func (app *application) contextSetSpyCat(r *http.Request, spyCat *model.SpyCat) *http.Request {
	ctx := context.WithValue(r.Context(), spyCatContextKey, spyCat)
	return r.WithContext(ctx)
}

func (app *application) contextGetSpyCat(r *http.Request) *model.SpyCat {
	user, ok := r.Context().Value(spyCatContextKey).(*model.SpyCat)
	if !ok {
		panic("missing user value in request context") //shrug
	}
	return user
}
