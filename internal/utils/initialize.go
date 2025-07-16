package utils

import (
	"blog-personal/internal/models"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Admin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"admin"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ConnectDatabase(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateAdminUser(db *sql.DB, cfg *Config) error {
	// Check if the admin user already exists
	var existingUser models.User
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", cfg.Admin.Username).Scan(&existingUser.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if existingUser.ID != 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (username, password, is_admin)
		VALUES (?, ?, ?)
	`
	_, err = db.Exec(query, cfg.Admin.Username, hashedPassword, true)
	return err
}

func ApplyMigrations(db *sql.DB, sqlPath string) error {
	data, err := os.ReadDir(sqlPath)
	if err != nil {
		return err
	}
	for _, file := range data {
		if file.IsDir() {
			continue
		}
		content, err := os.ReadFile(filepath.Join(sqlPath, file.Name()))
		if err != nil {
			return err
		}
		_, err = db.Exec(string(content))
		if err != nil {
			return err
		}
	}
	return nil
}
