package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/validator"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequestResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received := time.Now()
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		app.logger.Info("<-request:",
			"method", r.Method,
			"url", r.URL.String(),
			"ip", r.RemoteAddr,
			"body", string(bodyBytes),
			"received_at", received.Format(time.RFC3339),
		)

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		app.logger.Info("->response:",
			"code", rw.statusCode,
			"body", rw.body.String(),
			"took", time.Since(received),
		)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.contextSetSpyCat(r, model.AnonymousSpyCat)
			r = app.contextSetAgent(r, model.AnonymousAgent)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidCredentialsResponse(w, r)
			return
		}

		tokenPlaintext := headerParts[1]

		v := validator.New()

		if model.ValidateTokenPlaintext(v, tokenPlaintext); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token, err := app.tokensService.GetTokenByPlaintext(r.Context(), tokenPlaintext, model.ScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrorModelNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		switch token.UserType {
		case model.SpyCatUserType:
			spyCat, err := app.spyCatsService.GetById(r.Context(), token.UserID)
			if err != nil {
				switch {
				case errors.Is(err, storage.ErrorModelNotFound):
					app.invalidAuthenticationTokenResponse(w, r)
				default:
					app.serverErrorResponse(w, r, err)
				}
				return
			}
			r = app.contextSetSpyCat(r, spyCat)
			r = app.contextSetAgent(r, model.AnonymousAgent)
		case model.AgentUserType:
			agent, err := app.agentsService.GetById(r.Context(), token.UserID)
			if err != nil {
				switch {
				case errors.Is(err, storage.ErrorModelNotFound):
					app.invalidAuthenticationTokenResponse(w, r)
				default:
					app.serverErrorResponse(w, r, err)
				}
				return
			}
			r = app.contextSetAgent(r, agent)
			r = app.contextSetSpyCat(r, model.AnonymousSpyCat)
		default:
			panic(fmt.Sprintf("unsupported token type %s", token.UserType))
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireSpyCat(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spyCat := app.contextGetSpyCat(r)
		if spyCat.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAgent(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		agent := app.contextGetAgent(r)
		if agent.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
