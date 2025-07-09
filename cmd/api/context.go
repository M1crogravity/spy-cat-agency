package main

import (
	"context"
	"net/http"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
)

type contextKey string

const (
	spyCatContextKey = contextKey("spy-cat")
	agentContextKey  = contextKey("agent")
)

func (app *application) contextSetSpyCat(r *http.Request, spyCat *model.SpyCat) *http.Request {
	ctx := context.WithValue(r.Context(), spyCatContextKey, spyCat)
	return r.WithContext(ctx)
}

func (app *application) contextGetSpyCat(r *http.Request) *model.SpyCat {
	spyCat, ok := r.Context().Value(spyCatContextKey).(*model.SpyCat)
	if !ok {
		panic("missing spy cat value in request context")
	}
	return spyCat
}

func (app *application) contextSetAgent(r *http.Request, agent *model.Agent) *http.Request {
	ctx := context.WithValue(r.Context(), agentContextKey, agent)
	return r.WithContext(ctx)
}

func (app *application) contextGetAgent(r *http.Request) *model.Agent {
	agent, ok := r.Context().Value(agentContextKey).(*model.Agent)
	if !ok {
		panic("missing agent value in request context")
	}
	return agent
}
