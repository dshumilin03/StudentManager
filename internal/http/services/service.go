package services

import (
	"StudentManager/internal/domain"
	"StudentManager/internal/repository"
	"context"
	"log"
)

type StudentService interface {
	Create(ctx context.Context, fullName string, age int, groupNumber string, email string) (domain.Student, error)
	GetAll(ctx context.Context) ([]domain.Student, error)
}

type Services struct {
	Students StudentService
}

func NewServices(repo repository.StudentRepository) *Services {
	log.Printf("Services are created")
	return &Services{Students: NewStudentService(repo)}
}
