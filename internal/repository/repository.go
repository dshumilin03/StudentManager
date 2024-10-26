package repository

import (
	"StudentManager/internal/domain"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Repository[T any] interface {
	Create(ctx context.Context, domain T) (pgx.Rows, error)
	GetAll(ctx context.Context) (pgx.Rows, error)
	GetById(ctx context.Context, id int64) pgx.Row
	Update(ctx context.Context, domain T) (pgx.Rows, error)
	DeleteById(ctx context.Context, id int64) error
}

type StudentRepository interface {
	Repository[domain.Student]
	GetByEmail(ctx context.Context, email string) pgx.Row
}

type GroupRepository interface {
	Repository[domain.Group]
	GetByGroupNumber(ctx context.Context, name string) pgx.Row
}

type Repositories struct {
	Students StudentRepository
	Groups   GroupRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	log.Printf("Repositories are created")
	return &Repositories{
		Students: NewStudentRepositoryImplPostgres(db),
		Groups:   NewGroupRepositoryImplPostgres(db),
	}
}
