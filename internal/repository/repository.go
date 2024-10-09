package repository

import (
	"StudentManager/internal/domain"
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

type StudentRepository interface {
	Create(ctx context.Context, student domain.Student) error
	Get(ctx context.Context, id int) error
	Update(ctx context.Context, student domain.Student) error
	Delete(ctx context.Context, id int) error
}

type Repositories struct {
	Students StudentRepository
}

func NewRepositories(db *pgx.Conn) *Repositories {
	log.Printf("Repositories are created")
	return &Repositories{Students: NewStudentRepoPostgres(db)}
}