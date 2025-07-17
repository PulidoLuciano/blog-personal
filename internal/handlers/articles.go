package handlers

import (
	"blog-personal/internal/models"
	"database/sql"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

func ListArticles(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}
		page, er := strconv.Atoi(pageStr)
		if er != nil {
			http.Error(w, "Error al convertir la página a entero", http.StatusBadRequest)
			return
		}

		if page < 1 {
			page = 1
		}
		const pageSize = 10
		offset := (page - 1) * pageSize

		var articles []models.Article
		var rows *sql.Rows
		var err error

		if query != "" {
			rows, err = db.Query(`
			SELECT url, title, description, created_at, image 
			FROM articles 
			WHERE title LIKE ? OR description LIKE ? 
			ORDER BY created_at DESC 
			LIMIT ? OFFSET ?`, "%"+query+"%", "%"+query+"%", pageSize, offset)
		} else {
			rows, err = db.Query(`
			SELECT url, title, description, created_at, image 
			FROM articles 
			ORDER BY created_at DESC 
			LIMIT ? OFFSET ?`, pageSize, offset)
		}

		if err != nil {
			http.Error(w, "Error al obtener artículos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var a models.Article
			if err := rows.Scan(&a.Url, &a.Title, &a.Description, &a.CreatedAt, &a.Image); err == nil {
				articles = append(articles, a)
			}
		}

		// Para saber si hay más páginas
		var total int
		countQuery := "SELECT COUNT(*) FROM articles"
		if query != "" {
			countQuery = "SELECT COUNT(*) FROM articles WHERE title LIKE ? OR description LIKE ?"
			_ = db.QueryRow(countQuery, "%"+query+"%", "%"+query+"%").Scan(&total)
		} else {
			_ = db.QueryRow(countQuery).Scan(&total)
		}

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

		var info models.PersonalInfo
		er = db.QueryRow(`SELECT full_name, x, github, linkedin FROM personal_info LIMIT 1`).
			Scan(&info.FullName, &info.X, &info.GitHub, &info.LinkedIn)

		if er != nil {
			http.Error(w, "No se pudo obtener la información personal", http.StatusInternalServerError)
			return
		}

		data := struct {
			PageTitle    string
			Articles     []models.Article
			Query        string
			Page         int
			TotalPages   int
			PersonalInfo models.PersonalInfo
		}{
			PageTitle:    "Artículos",
			Articles:     articles,
			Query:        query,
			Page:         page,
			TotalPages:   totalPages,
			PersonalInfo: info,
		}

		var funcMap = template.FuncMap{
			"add": func(a, b int) int {
				return a + b
			},
			"sub": func(a, b int) int {
				return a - b
			},
		}

		tmpl, er := template.New("articles.html").Funcs(funcMap).ParseFiles("ui/html/articles.html")
		if er != nil {
			http.Error(w, "Error al parsear el template", http.StatusInternalServerError)
			log.Println(er)
			return
		}
		err = tmpl.ExecuteTemplate(w, "articles.html", data)
		if err != nil {
			http.Error(w, "Error al ejecutar el template", http.StatusInternalServerError)
			log.Printf("Error al ejecutar el template %v", err)
			return
		}
	}
}
