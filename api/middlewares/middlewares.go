package middlewares

import (
	"context"
	"net/http"

	"github.com/jameslahm/bloggy_backend/api/auth"
	"github.com/jameslahm/bloggy_backend/api/responses"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthtication(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, err := auth.ExtractTokenID(r)
		if err != nil {
			responses.ERROR(rw, http.StatusUnauthorized, err)
			return
		}
		ctx := context.WithValue(r.Context(), "id", id)
		r = r.WithContext(ctx)
		next(rw, r)
	}
}
