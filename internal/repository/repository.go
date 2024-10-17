package repository

import (
	"StudentManager/internal/domain"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// TODO make all repositories private
type StudentRepository interface {
	Create(ctx context.Context, student domain.Student) (pgx.Rows, error)
	GetById(ctx context.Context, id int64) pgx.Row
	Update(ctx context.Context, student domain.Student) (pgx.Rows, error)
	DeleteById(ctx context.Context, id int64) error
	GetAll(ctx context.Context) (pgx.Rows, error)
	GetByEmail(ctx context.Context, email string) pgx.Row
}

type GroupRepository interface {
	Create(ctx context.Context, group domain.Group) (pgx.Rows, error)
	GetById(ctx context.Context, id int64) pgx.Row
	Update(ctx context.Context, group domain.Group) (pgx.Rows, error)
	DeleteById(ctx context.Context, id int64) error
	GetAll(ctx context.Context) (pgx.Rows, error)
	GetByGroupNumber(ctx context.Context, name string) pgx.Row
}

type Repositories struct {
	Students StudentRepository
	Groups   GroupRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	log.Printf("Repositories are created")
	return &Repositories{
		Students: NewStudentRepoPostgres(db),
		Groups:   NewGroupRepoPostgres(db),
	}
}
