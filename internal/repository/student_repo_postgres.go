package repository

import (
	"StudentManager/internal/domain"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type StudentRepoPostgres struct {
	db *pgxpool.Pool
}

func NewStudentRepoPostgres(db *pgxpool.Pool) *StudentRepoPostgres {
	return &StudentRepoPostgres{
		db: db,
	}
}

func (repo *StudentRepoPostgres) GetAll(ctx context.Context) (pgx.Rows, error) {
	database := repo.db
	// TODO if returns null need to throw exception
	students, err := database.Query(ctx,
		"select * from student")
	if err != nil {
		log.Printf("%s: query executement", err)
		return nil, err
	}

	return students, err

}

func (repo *StudentRepoPostgres) Create(ctx context.Context, student domain.Student) (pgx.Rows, error) {
	database := repo.db

	_, err := database.Query(ctx,
		"insert into student(full_name, age, group_number, email) values($1, $2, $3, $4)",
		student.FullName, student.Age, student.GroupNumber, student.Email)
	if err != nil {
		log.Printf("%s: query executement", err)
		return nil, err
	}
	studentRow, err := database.Query(ctx, "select * from student where email = $1", student.Email)

	return studentRow, err
}
func (repo *StudentRepoPostgres) GetById(ctx context.Context, id int) pgx.Row {
	database := repo.db
	// TODO if returns null need to throw exception
	student := database.QueryRow(ctx,
		"select * from student where id = $1", id)

	return student
}
func (repo *StudentRepoPostgres) Update(ctx context.Context, student domain.Student) error {
	// TODO Implement
	database := repo.db
	_ = database
	return nil
}
func (repo *StudentRepoPostgres) Delete(ctx context.Context, id int) error {
	database := repo.db
	_, err := database.Exec(ctx, "delete from student where id = $1", id)
	if err != nil {
		fmt.Errorf("%s: query executement in deletion", err)
		return err
	}
	return err
}
