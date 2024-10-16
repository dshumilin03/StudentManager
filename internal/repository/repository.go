package repository

import (
	"StudentManager/internal/domain"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type StudentRepository interface {
	Create(ctx context.Context, student domain.Student) (pgx.Rows, error)
	GetById(ctx context.Context, id int64) pgx.Row
	Update(ctx context.Context, student domain.Student) (pgx.Rows, error)
	DeleteById(ctx context.Context, id int64) error
	GetAll(ctx context.Context) (pgx.Rows, error)
}

type Repositories struct {
	Students StudentRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	log.Printf("Repositories are created")
	return &Repositories{Students: NewStudentRepoPostgres(db)}
}
