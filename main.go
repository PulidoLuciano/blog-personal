package main

import (
	"blog-personal/internal/handlers"
	"blog-personal/internal/handlers/admin"
	"blog-personal/internal/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Cargar configuración
	cfg, err := utils.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Conectar a la base de datos
	db, err := utils.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Aplicar script SQL
	if err := utils.ApplyMigrations(db, "migrations/"); err != nil {
		log.Fatalf("Error al ejecutar las migraciones: %v", err)
	}

	// Crear usuario administrador
	if err := utils.CreateAdminUser(db, cfg); err != nil {
		log.Fatalf("Error al crear usuario administrador: %v", err)
	}

	db.Query(`
		INSERT INTO articles (url, title, description, content, image, author_id) VALUES (
			'articulo-prueba',
			"Articulo de prueba",
			"Descripción de prueba",
			"Contenido de prueba",
        '/static/assets/images/projects/sportsCalendar.webp',
        1
		);
	`)

	db.Query(`
		INSERT INTO articles (url, title, description, content, image, author_id) VALUES (
			'articulo-prueba-2',
			"Articulo de prueba 2",
			"Descripción de prueba 2",
			"Contenido de prueba 2",
        '/static/assets/images/projects/sportsCalendar.webp',
        1
		);
	`)

	db.Query(`
		INSERT INTO articles (url, title, description, content, image, author_id) VALUES (
			'articulo-prueba-3',
			"Articulo de prueba 3",
			"Descripción de prueba 3",
			"Contenido de prueba 3",
        '/static/assets/images/projects/sportsCalendar.webp',
        1
		);
	`)

	// Iniciar servidor HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.HomeHandler(db))
	router.HandleFunc("/admin", admin.DashboardHandler)
	router.HandleFunc("/admin/personal", admin.PersonalInfoGetHandler(db)).Methods("GET")
	router.HandleFunc("/admin/personal", admin.PersonalInfoPostHandler(db)).Methods("POST")

	fs := http.FileServer(http.Dir("ui/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Error al iniciar servidor HTTP: %v", err)
	}
	fmt.Printf("Servidor HTTP iniciado en http://%s\n", addr)
}
