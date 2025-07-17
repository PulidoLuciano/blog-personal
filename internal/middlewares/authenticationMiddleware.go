package middlewares

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AuthenticationMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session")
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
