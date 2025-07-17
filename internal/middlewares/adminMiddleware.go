package middlewares

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AdminMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session")
			if admin, ok := session.Values["is_admin"].(bool); !ok || !admin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
