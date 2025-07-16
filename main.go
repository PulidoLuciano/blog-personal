package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Reemplazá estos valores por los correctos
	dsn := "webuser:webpassword@tcp(172.16.90.131:3306)/blog"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al abrir conexión: %s", err)
	}

	// Ping con timeout de red
	err = db.Ping()
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %s", err)
	}

	fmt.Println("Conexión exitosa a MariaDB")
}