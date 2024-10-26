package service

import (
	"StudentManager/internal/domain"
	"StudentManager/internal/dto"
	"StudentManager/internal/repository"
	"context"
	"log"
)

type Service[Dto any, Domain any] interface {
	Create(ctx context.Context, dto Dto) (Domain, error)
	GetAll(ctx context.Context) ([]Domain, error)
	GetById(ctx context.Context, id int64) (Domain, error)
	Update(ctx context.Context, dto Dto) (Domain, error)
	DeleteById(ctx context.Context, id int64) error
	GetService() Service[Dto, Domain]
}

type StudentService interface {
	Service[dto.StudentDto, domain.Student]
	IsStudentExistsByEmail(ctx context.Context, email string) bool
	IsStudentExistsById(ctx context.Context, id int64) bool
}

type GroupService interface {
	Service[dto.GroupDto, domain.Group]
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
