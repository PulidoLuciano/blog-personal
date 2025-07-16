package admin

import (
	"blog-personal/internal/models"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var initialPath = "ui/"
var folderProjects = "static/assets/images/projects/"

func ProjectsList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query(`SELECT id, title, description, url_repository, url_website, image, is_hidden FROM projects`)
		if err != nil {
			http.Error(w, "Error al cargar proyectos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var projects []models.Project
		for rows.Next() {
			var p models.Project
			if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.UrlRepository, &p.UrlWebsite, &p.Image, &p.IsHidden); err == nil {
				projects = append(projects, p)
			}
		}

		tmpl := template.Must(template.ParseFiles("ui/html/admin/projects.html"))
		tmpl.Execute(w, projects)
	}
}

func ProjectForm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		var p models.Project

		if idStr != "" {
			id, _ := strconv.Atoi(idStr)
			err := db.QueryRow(`SELECT id, title, description, url_repository, url_website, image, is_hidden FROM projects WHERE id = ?`, id).
				Scan(&p.ID, &p.Title, &p.Description, &p.UrlRepository, &p.UrlWebsite, &p.Image, &p.IsHidden)
			if err != nil {
				http.Error(w, "Proyecto no encontrado", http.StatusNotFound)
				return
			}
		}

		tmpl := template.Must(template.ParseFiles("ui/html/admin/projects_form.html"))
		tmpl.Execute(w, p)
	}
}

func ProjectSave(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, "Error al convertir el ID", http.StatusBadRequest)
			log.Println("Error al convertir el ID:", err)
			return
		}

		p := models.Project{
			Title:         r.FormValue("title"),
			Description:   r.FormValue("description"),
			UrlRepository: r.FormValue("url_repository"),
			UrlWebsite:    r.FormValue("url_website"),
			IsHidden:      r.FormValue("is_hidden") == "on",
			Image:         r.FormValue("image-path"),
		}

		file, header, err := r.FormFile("image")
		if err == nil && header.Filename != "" {
			defer file.Close()

			imageName := fmt.Sprintf("%d%s", id, filepath.Ext(header.Filename))
			filePath := filepath.Join(folderProjects, imageName)
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
			p.Image = filePath
		}

		if id > 0 {
			_, err := db.Exec(`UPDATE projects SET title=?, description=?, url_repository=?, url_website=?, image=?, is_hidden=? WHERE id=?`,
				p.Title, p.Description, p.UrlRepository, p.UrlWebsite, p.Image, p.IsHidden, id)
			if err != nil {
				http.Error(w, "Error al actualizar", http.StatusInternalServerError)
				return
			}
		} else {
			_, err := db.Exec(`INSERT INTO projects (title, description, url_repository, url_website, image, is_hidden) VALUES (?, ?, ?, ?, ?, ?)`,
				p.Title, p.Description, p.UrlRepository, p.UrlWebsite, p.Image, p.IsHidden)
			if err != nil {
				http.Error(w, "Error al crear", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/admin/projects", http.StatusSeeOther)
	}
}

func ProjectDelete(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("id"))

		var imagePath string
		err := db.QueryRow(`SELECT image FROM projects WHERE id = ?`, id).Scan(&imagePath)
		if err != nil {
			http.Error(w, "Error al obtener la imagen", http.StatusInternalServerError)
			return
		}
		os.Remove(filepath.Join(initialPath, imagePath))

		_, err = db.Exec(`DELETE FROM projects WHERE id = ?`, id)
		if err != nil {
			http.Error(w, "Error al eliminar", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/projects", http.StatusSeeOther)
	}
}
