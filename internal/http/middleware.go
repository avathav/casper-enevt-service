package http

import (
	"net/http"

	internalgorm "event-service/internal/database/gorm"

	"gorm.io/gorm"
)

func GORMConnectionMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := internalgorm.ContextWithConnection(r.Context(), db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
