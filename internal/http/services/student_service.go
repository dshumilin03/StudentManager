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
	studentRepository repository.StudentRepository
	groupService      GroupService
}

func NewStudentServiceImpl(repo repository.StudentRepository, service GroupService) *StudentServiceImpl {
	return &StudentServiceImpl{
		studentRepository: repo,
		groupService:      service,
	}
}

func (studentService *StudentServiceImpl) Create(
	ctx context.Context,
	fullName string,
	age int,
	groupNumber string,
	email string,
) (domain.Student, error) {

	repo := studentService.studentRepository
	groupService := studentService.groupService

	student := domain.Student{
		FullName:    fullName,
		Age:         age,
		GroupNumber: groupNumber,
		Email:       email,
	}

	if studentService.IsStudentExistsByEmail(ctx, email) {
		log.Println("student already exists")
		return domain.Student{}, errors.New("student already exists")
	}

	if !groupService.IsGroupExistsByNumber(ctx, groupNumber) {
		log.Println("group doesn't exist")
		return domain.Student{}, errors.New("group doesn't exist")
	}

	studentRow, err := repo.Create(ctx, student)
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

func (studentService *StudentServiceImpl) GetAll(ctx context.Context) ([]domain.Student, error) {
	service := studentService.studentRepository

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

func (studentService *StudentServiceImpl) GetById(ctx context.Context, id int64) (domain.Student, error) {
	service := studentService.studentRepository

	row := service.GetById(ctx, id)

	students, err := convertStudentRowToDomain(row)
	if err != nil {
		log.Printf("failed to convert students into domain %v", err)

	}

	return students, nil
}

func (studentService *StudentServiceImpl) Update(
	ctx context.Context,
	id int64,
	fullName string,
	age int,
	groupNumber string,
	email string,
) (domain.Student, error) {
	repo := studentService.studentRepository

	student := domain.Student{
		Id:          id,
		FullName:    fullName,
		Age:         age,
		GroupNumber: groupNumber,
		Email:       email,
	}

	if !studentService.IsStudentExistsById(ctx, id) {
		log.Println("student doesn't exists")
		return domain.Student{}, errors.New("student doesn't exists")
	}

	studentRow, err := repo.Update(ctx, student)
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

func (studentService *StudentServiceImpl) DeleteById(ctx context.Context, id int64) error {
	repo := studentService.studentRepository

	if !studentService.IsStudentExistsById(ctx, id) {
		log.Println("student doesn't exists")
		return errors.New("student does not exist")
	}
	err := repo.DeleteById(ctx, id)
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

func (studentService *StudentServiceImpl) IsStudentExistsByEmail(ctx context.Context, email string) bool {
	service := studentService.studentRepository

	if errors.Is(service.GetByEmail(ctx, email).Scan(), pgx.ErrNoRows) {
		return false
	}

	return true
}

func (studentService *StudentServiceImpl) IsStudentExistsById(ctx context.Context, id int64) bool {
	service := studentService.studentRepository

	if errors.Is(service.GetById(ctx, id).Scan(), pgx.ErrNoRows) {
		return false
	}

	return true
}
