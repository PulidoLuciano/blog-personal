package handlers

import (
	"blog-personal/internal/models"
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type HomeData struct {
	PageTitle    string
	PersonalInfo models.PersonalInfo
	Articles     []models.Article
	Projects     []models.Project
}

func HomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Obtener información personal
		rows, err := db.Query(`
			SELECT full_name, description, bio, profile_image, x, github, linkedin
			FROM personal_info
			LIMIT 1`)
		if err != nil {
			http.Error(w, "Error al obtener información personal", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()

		var personalInfo models.PersonalInfo
		if rows.Next() {
			var i models.PersonalInfo
			if err := rows.Scan(&i.FullName, &i.Description, &i.Bio, &i.ProfileImage, &i.X,
				&i.GitHub, &i.LinkedIn); err == nil {
				personalInfo = i
			}
		}

		// Obtener artículos
		rows, err = db.Query(`
			SELECT url, title, created_at, image
			FROM articles 
			ORDER BY created_at DESC 
			LIMIT 3`)
		if err != nil {
			http.Error(w, "Error al obtener artículos", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()

		var articles []models.Article
		for rows.Next() {
			var a models.Article
			if err := rows.Scan(&a.Url, &a.Title, &a.CreatedAt, &a.Image); err == nil {
				articles = append(articles, a)
			}
		}

		// Obtener proyectos
		rows, err = db.Query(`
			SELECT title, description, url_repository, url_website, image
			FROM projects 
			WHERE is_hidden = FALSE`)
		if err != nil {
			http.Error(w, "Error al obtener proyectos", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()

		var projects []models.Project
		for rows.Next() {
			var p models.Project
			if err := rows.Scan(&p.Title, &p.Description, &p.UrlRepository, &p.UrlWebsite, &p.Image); err == nil {
				projects = append(projects, p)
			}
		}

		tmpl := template.Must(template.ParseFiles("ui/html/index.html"))

		data := HomeData{
			PageTitle:    "Blog de Luciano Pulido",
			PersonalInfo: personalInfo,
			Articles:     articles,
			Projects:     projects,
		}

		tmpl.Execute(w, data)
	}
}
