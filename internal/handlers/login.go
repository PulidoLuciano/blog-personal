package handlers

import (
	"blog-personal/internal/models"
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func LoginForm(store *sessions.CookieStore, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		if auth, ok := session.Values["authenticated"].(bool); ok && auth {
			referer := r.Header.Get("Referer")
			if referer == "" {
				referer = "/"
			}
			http.Redirect(w, r, referer, http.StatusSeeOther)
			return
		}

		var info models.PersonalInfo
		err := db.QueryRow(`SELECT id, full_name, description, bio, profile_image, x, github, linkedin FROM personal_info LIMIT 1`).
			Scan(&info.ID, &info.FullName, &info.Description, &info.Bio, &info.ProfileImage, &info.X, &info.GitHub, &info.LinkedIn)

		if err != nil {
			http.Error(w, "No se pudo obtener la información personal", http.StatusInternalServerError)
			return
		}

		flashSession, _ := store.Get(r, "flash")

		fm := flashSession.Flashes("login-error")
		data := map[string]interface{}{
			"PersonalInfo": info,
		}
		if fm != nil {
			data["Error"] = fm[0].(string)
		}

		tmpl := template.Must(template.ParseFiles("ui/html/login.html"))
		tmpl.Execute(w, data)
	}
}

func LogoutHandler(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		session.Values["authenticated"] = false
		session.Save(r, w)

		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "/"
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
	}
}

func LoginHandler(store *sessions.CookieStore, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		flashSession, _ := store.Get(r, "flash")

		info := models.User{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		rows, err := db.Query(`SELECT id, username, password, is_admin FROM users WHERE username = ?`, info.Username)
		if err != nil {
			flashSession.AddFlash("Usuario o contraseña incorrectos", "login-error")
			flashSession.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		defer rows.Close()

		var user models.User
		for rows.Next() {
			err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin)
			if err != nil {
				flashSession.AddFlash("Error al obtener usuario", "login-error")
				flashSession.Save(r, w)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password))
		if err != nil {
			flashSession.AddFlash("Usuario o contraseña incorrectos", "login-error")
			flashSession.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, _ := store.Get(r, "session")

		session.Values["authenticated"] = true
		session.Values["user_id"] = user.ID
		session.Values["username"] = user.Username
		session.Values["is_admin"] = user.IsAdmin
		session.Save(r, w)

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}
