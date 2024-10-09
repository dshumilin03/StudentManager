package repository

import (
	"StudentManager/internal/domain"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type StudentRepoPostgres struct {
	db *pgx.Conn
}

func NewStudentRepoPostgres(db *pgx.Conn) *StudentRepoPostgres {
	return &StudentRepoPostgres{
		db: db,
	}
}

func (repo *StudentRepoPostgres) Create(ctx context.Context, student domain.Student) error {
	database := repo.db

	_, err := database.Exec(ctx,
		"insert into student(full_name, age, group_number, email) values($1, $2, $3, $4)",
		student.FullName, student.Age, student.GroupNumber, student.Email)
	if err != nil {
		fmt.Errorf("%s: query executement", err)
		return err
	}
	return err
}
func (repo *StudentRepoPostgres) Get(ctx context.Context, id int) error {
	// TODO Implement
	database := repo.db
	_ = database
	return nil
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
