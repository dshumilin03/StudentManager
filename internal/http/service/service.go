package service

import (
	"StudentManager/internal/domain"
	"StudentManager/internal/dto"
	"StudentManager/internal/repository"
	"context"
	"log"
)

type StudentService interface {
	Create(ctx context.Context, dto dto.StudentDto) (domain.Student, error)
	GetAll(ctx context.Context) ([]domain.Student, error)
	GetById(ctx context.Context, id int64) (domain.Student, error)
	Update(ctx context.Context, dto dto.StudentDto) (domain.Student, error)
	DeleteById(ctx context.Context, id int64) error
	IsStudentExistsByEmail(ctx context.Context, email string) bool
	IsStudentExistsById(ctx context.Context, id int64) bool
}

type GroupService interface {
	Create(ctx context.Context, dto dto.GroupDto) (domain.Group, error)
	GetAll(ctx context.Context) ([]domain.Group, error)
	GetById(ctx context.Context, id int64) (domain.Group, error)
	Update(ctx context.Context, dto dto.GroupDto) (domain.Group, error)
	DeleteById(ctx context.Context, id int64) error
	IsGroupExistsByNumber(ctx context.Context, groupNumber string) bool
	IsGroupExistsById(ctx context.Context, id int64) bool
}

type Services struct {
	Students StudentService
	Groups   GroupService
}

func NewServices(repositories *repository.Repositories) *Services {
	log.Printf("Services are created")
	return &Services{
		Students: NewStudentServiceImpl(repositories.Students, NewGroupServiceImpl(repositories.Groups)),
		Groups:   NewGroupServiceImpl(repositories.Groups),
	}
}
