package services

import (
	"StudentManager/internal/domain"
	"StudentManager/internal/repository"
	"context"
	"log"
)

// StudentService TODO try to put arguments for student domain into one type
type StudentService interface {
	Create(ctx context.Context, fullName string, age int, groupNumber string, email string) (domain.Student, error)
	GetAll(ctx context.Context) ([]domain.Student, error)
	GetById(ctx context.Context, id int64) (domain.Student, error)
	Update(ctx context.Context, id int64, fullName string, age int, groupNumber string, email string) (domain.Student, error)
	DeleteById(ctx context.Context, id int64) error
}

type Services struct {
	Students StudentService
}

func NewServices(repo repository.StudentRepository) *Services {
	log.Printf("Services are created")
	return &Services{Students: NewStudentService(repo)}
}
