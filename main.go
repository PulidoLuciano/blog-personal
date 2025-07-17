package main

import (
	"blog-personal/internal/handlers"
	"blog-personal/internal/handlers/admin"
	"blog-personal/internal/middlewares"
	"blog-personal/internal/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	// Cargar configuración
	cfg, err := utils.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	key := []byte(cfg.SecretKey)
	store := sessions.NewCookieStore(key)

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

	// Iniciar servidor HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.HomeHandler(db))
	router.HandleFunc("/login", handlers.LoginHandler(store, db)).Methods("POST")
	router.HandleFunc("/login", handlers.LoginForm(store, db)).Methods("GET")
	router.HandleFunc("/logout", handlers.LogoutHandler(store)).Methods("POST")
	router.HandleFunc("/articles", handlers.ListArticles(db)).Methods("GET")
	router.HandleFunc("/articles/{article-id}", handlers.Article(db)).Methods("GET")

	// Admin
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewares.AuthenticationMiddleware(store), middlewares.AdminMiddleware(store))

	adminRouter.HandleFunc("", admin.DashboardHandler)
	// Personal
	adminRouter.HandleFunc("/personal", admin.PersonalInfoGetHandler(db)).Methods("GET")
	adminRouter.HandleFunc("/personal", admin.PersonalInfoPostHandler(db)).Methods("POST")
	// Projects
	adminRouter.HandleFunc("/projects", admin.ProjectsList(db)).Methods("GET")
	adminRouter.HandleFunc("/projects/new", admin.ProjectForm(db)).Methods("GET")
	adminRouter.HandleFunc("/projects/new", admin.ProjectSave(db)).Methods("POST")
	adminRouter.HandleFunc("/projects/edit", admin.ProjectForm(db)).Methods("GET")
	adminRouter.HandleFunc("/projects/edit", admin.ProjectSave(db)).Methods("POST")
	adminRouter.HandleFunc("/projects/delete", admin.ProjectDelete(db)).Methods("POST")
	// Articles
	adminRouter.HandleFunc("/articles", admin.ArticlesList(db)).Methods("GET")
	adminRouter.HandleFunc("/articles/new", admin.ArticleForm(db)).Methods("GET")
	adminRouter.HandleFunc("/articles/new", admin.ArticleSave(db)).Methods("POST")
	adminRouter.HandleFunc("/articles/edit", admin.ArticleForm(db)).Methods("GET")
	adminRouter.HandleFunc("/articles/edit", admin.ArticleSave(db)).Methods("POST")
	adminRouter.HandleFunc("/articles/delete", admin.ArticleDelete(db)).Methods("POST")

	// Authentication

	fs := http.FileServer(http.Dir("ui/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Error al iniciar servidor HTTP: %v", err)
	}
	fmt.Printf("Servidor HTTP iniciado en http://%s\n", addr)
}
