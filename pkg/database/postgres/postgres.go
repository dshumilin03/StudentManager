package postgres

import (
	"StudentManager/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	db *sql.DB
}

type Database struct {
	Username string
	Password string
	Host     string
	Port     string
	Dbname   string
	Sslmode  string
}

func New(database config.Database) (*Storage, error) {

	// Формирование строки подключения
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		database.Username,
		database.Password,
		database.Host,
		database.Port,
		database.Dbname,
		database.Sslmode)

	db, err := sql.Open("postgres", connStr)
	// Подключение к базе данных
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to the database!")
	return &Storage{db: db}, nil
}
