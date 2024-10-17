package services

import (
	"StudentManager/internal/domain"
	"StudentManager/internal/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"log"
)

type GetStudentRequest struct {
	FullName string `json:"full_name" env-required:"true"`
}

type StudentServiceImpl struct {
	repo repository.StudentRepository
}

func NewStudentServiceImpl(repo repository.StudentRepository) *StudentServiceImpl {
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

	if repo.isStudentExistsByEmail(ctx, email) {
		log.Println("student already exists")
		return domain.Student{}, errors.New("student already exists")
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

func (repo *StudentServiceImpl) GetById(ctx context.Context, id int64) (domain.Student, error) {
	service := repo.repo

	row := service.GetById(ctx, id)

	students, err := convertStudentRowToDomain(row)
	if err != nil {
		log.Printf("failed to convert students into domain %v", err)

	}

	return students, nil
}

func (repo *StudentServiceImpl) Update(
	ctx context.Context,
	id int64,
	fullName string,
	age int,
	groupNumber string,
	email string,
) (domain.Student, error) {
	service := repo.repo

	student := domain.Student{
		Id:          id,
		FullName:    fullName,
		Age:         age,
		GroupNumber: groupNumber,
		Email:       email,
	}

	if !repo.isStudentExistsById(ctx, id) {
		log.Println("student doesn't exists")
		return domain.Student{}, errors.New("student doesn't exists")
	}

	studentRow, err := service.Update(ctx, student)
	if err != nil {
		log.Printf("failed to update student %v", err)
		return domain.Student{}, err
	}

	updatedStudent, err := convertStudentsRowsToDomain(studentRow)
	if err != nil {
		log.Printf("failed to convert student into domain %v", err)
		return domain.Student{}, err
	}

	return updatedStudent[0], err
}

func (repo *StudentServiceImpl) DeleteById(ctx context.Context, id int64) error {
	service := repo.repo

	if !repo.isStudentExistsById(ctx, id) {
		log.Println("student doesn't exists")
		return errors.New("student does not exist")
	}
	err := service.DeleteById(ctx, id)
	if err != nil {
		log.Printf("failed to delete student %v", err)
	}

	return nil
}

func convertStudentRowToDomain(row pgx.Row) (domain.Student, error) {
	var student domain.Student

	err := row.Scan(&student.Id, &student.FullName, &student.Age, &student.GroupNumber, &student.Email)

	if err != nil {
		return domain.Student{}, err
	}

	return student, err
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

func (repo *StudentServiceImpl) isStudentExistsByEmail(ctx context.Context, email string) bool {
	service := repo.repo

	if errors.Is(service.GetByEmail(ctx, email).Scan(), pgx.ErrNoRows) {
		return false
	}

	return true
}

func (repo *StudentServiceImpl) isStudentExistsById(ctx context.Context, id int64) bool {
	service := repo.repo

	if errors.Is(service.GetById(ctx, id).Scan(), pgx.ErrNoRows) {
		return false
	}

	return true
}
