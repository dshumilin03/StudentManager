package postgres

import (
	"StudentManager/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func New(database config.Database) (*pgxpool.Pool, error) {

	// Формирование строки подключения
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		database.Username,
		database.Password,
		database.Host,
		database.Port,
		database.Dbname,
		database.Sslmode)

	_, err := pgx.ParseConfig(connStr)
	if err != nil {
		log.Fatal(err)
	}

	client, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	// Проверка соединения
	if err := client.Ping(context.Background()); err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to the database!")
	return client, nil
}
