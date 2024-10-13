package services

import (
	"StudentManager/internal/domain"
	resp "StudentManager/internal/http/response"
	"StudentManager/internal/repository"
	"context"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"log"
	"net/http"
)

type GetStudentRequest struct {
	FullName string `json:"full_name" env-required:"true"`
}

type Response struct {
	resp.Response
}

type StudentServiceImpl struct {
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) *StudentServiceImpl {
	return &StudentServiceImpl{
		repo: repo,
	}
}

func (repo *StudentServiceImpl) Create(
	ctx context.Context,
	fullName string,
	age int,
	groupNumber string,
	email string,
) (domain.Student, error) {
	service := repo.repo

	student := domain.Student{
		FullName:    fullName,
		Age:         age,
		GroupNumber: groupNumber,
		Email:       email,
	}

	studentRow, err := service.Create(ctx, student)
	if err != nil {
		log.Printf("failed to create student %v", err)
		return domain.Student{}, err
	}

	createdStudent, err := convertStudentsRowsToDomain(studentRow)
	if err != nil {
		log.Printf("failed to convert student into domain %v", err)
		return domain.Student{}, err
	}

	return createdStudent[0], err
}

func (repo *StudentServiceImpl) GetAll(ctx context.Context) ([]domain.Student, error) {
	service := repo.repo

	rows, err := service.GetAll(ctx)
	if err != nil {
		log.Printf("failed to get students %v", err)
	}

	students, err := convertStudentsRowsToDomain(rows)
	if err != nil {
		log.Printf("failed to convert students into domain %v", err)

	}

	return students, nil
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}

func convertStudentsRowsToDomain(rows pgx.Rows) ([]domain.Student, error) {
	var students []domain.Student

	for rows.Next() {
		var r domain.Student
		err := rows.Scan(&r.Id, &r.FullName, &r.Age, &r.GroupNumber, &r.Email)
		if err != nil {
			return nil, err
		}
		students = append(students, r)
	}

	return students, nil
}
