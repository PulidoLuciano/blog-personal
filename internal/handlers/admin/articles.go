package admin

import (
	"blog-personal/internal/models"
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var folderArticles = "static/assets/images/articles/"

func ArticlesList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query(`SELECT url, title, created_at FROM articles`)
		if err != nil {
			http.Error(w, "Error al cargar artículos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var articles []models.Article
		for rows.Next() {
			var a models.Article
			if err := rows.Scan(&a.Url, &a.Title, &a.CreatedAt); err == nil {
				articles = append(articles, a)
			}
		}

		tmpl := template.Must(template.ParseFiles("ui/html/admin/articles.html"))
		tmpl.Execute(w, articles)
	}
}

func ArticleForm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlStr := r.URL.Query().Get("url")
		var a models.Article

		if urlStr != "" {
			err := db.QueryRow(`SELECT url, title, description, content, image FROM articles WHERE url = ?`, urlStr).
				Scan(&a.Url, &a.Title, &a.Description, &a.Content, &a.Image)
			if err != nil {
				http.Error(w, "Artículo no encontrado", http.StatusNotFound)
				return
			}
		}

		tmpl := template.Must(template.ParseFiles("ui/html/admin/articles_form.html"))
		tmpl.Execute(w, a)
	}
}

func ArticleSave(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		currentUrl := r.FormValue("url-actual")
		a := models.Article{
			Url:         r.FormValue("url"),
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Content:     r.FormValue("content"),
			Image:       r.FormValue("image-path"),
			AuthorId:    1,
		}

		file, header, err := r.FormFile("image")
		if err == nil && header.Filename != "" {
			defer file.Close()

			imageName := a.Url
			filePath := filepath.Join(folderArticles, imageName)
			savePath := filepath.Join(initialPath, filePath)

			// Crear archivo en disco
			dst, err := os.Create(savePath)
			if err != nil {
				http.Error(w, "No se pudo guardar la imagen", http.StatusInternalServerError)
				log.Println("Error al crear el archivo:", err)
				return
			}
			defer dst.Close()

			// Copiar contenido
			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, "Error al guardar imagen", http.StatusInternalServerError)
				log.Println("Error al guardar imagen:", err)
				return
			}
			a.Image = filePath
		}

		if currentUrl != "" && a.Url == currentUrl {
			_, err := db.Exec(`UPDATE articles SET title=?, description=?, content=?, image=? WHERE url=?`,
				a.Title, a.Description, a.Content, a.Image, a.Url)
			if err != nil {
				http.Error(w, "Error al actualizar", http.StatusInternalServerError)
				return
			}
		} else {
			if currentUrl != "" {
				os.Remove(filepath.Join(initialPath, folderArticles, currentUrl))
				_, err := db.Exec(`DELETE FROM articles WHERE url = ?`, currentUrl)
				if err != nil {
					http.Error(w, "Error al eliminar", http.StatusInternalServerError)
					return
				}
			}
			_, err := db.Exec(`INSERT INTO articles (url, title, description, content, image, author_id) VALUES (?, ?, ?, ?, ?, ?)`,
				a.Url, a.Title, a.Description, a.Content, a.Image, a.AuthorId)
			if err != nil {
				http.Error(w, "Error al crear", http.StatusInternalServerError)
				log.Println("Error al crear:", err)
				return
			}
		}

		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}
}

func ArticleDelete(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		url := r.FormValue("url")

		var imagePath string
		err := db.QueryRow(`SELECT image FROM articles WHERE url = ?`, url).Scan(&imagePath)
		if err != nil {
			http.Error(w, "Error al obtener la imagen", http.StatusInternalServerError)
			return
		}
		os.Remove(filepath.Join(initialPath, imagePath))

		_, err = db.Exec(`DELETE FROM articles WHERE url = ?`, url)
		if err != nil {
			http.Error(w, "Error al eliminar", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}
}
