package repository

import (
	"StudentManager/internal/domain"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type GroupRepoPostgres struct {
	db *pgxpool.Pool
}

func NewGroupRepoPostgres(db *pgxpool.Pool) *GroupRepoPostgres {
	return &GroupRepoPostgres{
		db: db,
	}
}

func (repo *GroupRepoPostgres) GetAll(ctx context.Context) (pgx.Rows, error) {
	database := repo.db

	groups, err := database.Query(ctx,
		"select * from \"group\" order by id")
	if err != nil {
		log.Printf("%s: query executement", err)
		return nil, err
	}

	return groups, err

}

func (repo *GroupRepoPostgres) Create(ctx context.Context, group domain.Group) (pgx.Rows, error) {
	database := repo.db

	_, err := database.Query(ctx,
		"insert into \"group\"(group_number) values($1)",
		group.GroupNumber)
	if err != nil {
		log.Printf("%s: query executement", err)
		return nil, err
	}
	groupRows, err := database.Query(ctx, "select * from \"group\" where group_number = $1", group.GroupNumber)

	return groupRows, err
}
func (repo *GroupRepoPostgres) GetById(ctx context.Context, id int64) pgx.Row {
	database := repo.db

	group := database.QueryRow(ctx,
		"select * from \"group\" where id = $1", id)

	return group
}
func (repo *GroupRepoPostgres) Update(ctx context.Context, group domain.Group) (pgx.Rows, error) {
	database := repo.db

	_, err := database.Query(ctx,
		"update \"group\" set group_number = $1", group.GroupNumber)
	if err != nil {
		log.Printf("%s: query executement or group doesn't exists", err)
		return nil, err
	}
	groupRows, err := database.Query(ctx, "select * from \"group\" where id = $1", group.Id)

	return groupRows, err
}
func (repo *GroupRepoPostgres) DeleteById(ctx context.Context, id int64) error {
	database := repo.db
	_, err := database.Exec(ctx, "delete from \"group\" where id = $1", id)
	if err != nil {
		fmt.Errorf("%s: query executement in deletion", err)
		return err
	}
	return err
}

func (repo *GroupRepoPostgres) GetByGroupNumber(ctx context.Context, groupNumber string) pgx.Row {
	database := repo.db

	group := database.QueryRow(ctx,
		"select * from \"group\" where group_number = $1", groupNumber)

	return group
}
