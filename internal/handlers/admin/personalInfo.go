package admin

import (
	"blog-personal/internal/models"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func PersonalInfoGetHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var info models.PersonalInfo
		err := db.QueryRow(`SELECT id, full_name, description, bio, profile_image, x, github, linkedin FROM personal_info LIMIT 1`).
			Scan(&info.ID, &info.FullName, &info.Description, &info.Bio, &info.ProfileImage, &info.X, &info.GitHub, &info.LinkedIn)

		if err != nil {
			http.Error(w, "No se pudo obtener la información personal", http.StatusInternalServerError)
			log.Println("DB error:", err)
			return
		}

		tmpl, err := template.ParseFiles("ui/html/admin/personal_info.html")
		if err != nil {
			http.Error(w, "Error cargando plantilla", http.StatusInternalServerError)
			log.Println("Template error:", err)
			return
		}

		tmpl.Execute(w, info)
	}
}

func PersonalInfoPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		info := models.PersonalInfo{
			ID:          id,
			FullName:    r.FormValue("full_name"),
			Description: r.FormValue("description"),
			Bio:         r.FormValue("bio"),
			X:           r.FormValue("x"),
			GitHub:      r.FormValue("github"),
			LinkedIn:    r.FormValue("linkedin"),
		}

		_, err = db.Exec(`
			UPDATE personal_info SET
				full_name = ?,
				description = ?,
				bio = ?,
				x = ?,
				github = ?,
				linkedin = ?
			WHERE id = ?`,
			info.FullName, info.Description, info.Bio, info.X, info.GitHub, info.LinkedIn, info.ID)

		if err != nil {
			http.Error(w, "No se pudo actualizar la información", http.StatusInternalServerError)
			log.Println("Update error:", err)
			return
		}

		http.Redirect(w, r, "/admin/personal", http.StatusSeeOther)
	}
}
