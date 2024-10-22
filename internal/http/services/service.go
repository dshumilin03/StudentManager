package services

import (
	"StudentManager/internal/domain"
	"StudentManager/internal/repository"
	"context"
	"log"
)

// StudentService TODO try to put arguments for student domain into one type
// TODO make all services private
type StudentService interface {
	Create(ctx context.Context, fullName string, age int, groupNumber string, email string) (domain.Student, error)
	GetAll(ctx context.Context) ([]domain.Student, error)
	GetById(ctx context.Context, id int64) (domain.Student, error)
	Update(ctx context.Context, id int64, fullName string, age int, groupNumber string, email string) (domain.Student, error)
	DeleteById(ctx context.Context, id int64) error
	IsStudentExistsByEmail(ctx context.Context, email string) bool
	IsStudentExistsById(ctx context.Context, id int64) bool
}

type GroupService interface {
	Create(ctx context.Context, groupNumber string) (domain.Group, error)
	GetAll(ctx context.Context) ([]domain.Group, error)
	GetById(ctx context.Context, id int64) (domain.Group, error)
	Update(ctx context.Context, id int64, groupNumber string) (domain.Group, error)
	DeleteById(ctx context.Context, id int64) error
	IsGroupExistsByNumber(ctx context.Context, groupNumber string) bool
	IsGroupExistsById(ctx context.Context, id int64) bool
}

type Services struct {
	Students StudentService
	Groups   GroupService
}

// NewServices TODO get all repos in one
func NewServices(studentRepo repository.StudentRepository, groupRepo repository.GroupRepository) *Services {
	log.Printf("Services are created")
	return &Services{
		Students: NewStudentServiceImpl(studentRepo, NewGroupServiceImpl(groupRepo)),
		Groups:   NewGroupServiceImpl(groupRepo),
	}
}
