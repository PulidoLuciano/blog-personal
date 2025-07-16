package main

import (
	"blog-personal/internal/utils"
	"fmt"
	"log"
	"net/http"
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

	// Iniciar servidor HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "¡Blog personal corriendo en Go y MySQL!")
	})

	fs := http.FileServer(http.Dir("ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error al iniciar servidor HTTP: %v", err)
	}
	fmt.Printf("Servidor HTTP iniciado en http://%s\n", addr)
}
