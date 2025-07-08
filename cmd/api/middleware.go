package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
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

		app.logger.Info("request:",
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

		app.logger.Info("response:",
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
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
