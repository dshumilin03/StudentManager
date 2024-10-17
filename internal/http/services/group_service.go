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

type GetGroupRequest struct {
	GroupNumber string `json:"group_number" env-required:"true"`
}

type GroupServiceImpl struct {
	repo repository.GroupRepository
}

func NewGroupServiceImpl(repo repository.GroupRepository) *GroupServiceImpl {
	return &GroupServiceImpl{
		repo: repo,
	}
}

func (repo *GroupServiceImpl) Create(
	ctx context.Context, groupNumber string) (domain.Group, error) {
	service := repo.repo

	group := domain.Group{
		GroupNumber: groupNumber,
	}

	if repo.isGroupExistsByNumber(ctx, groupNumber) {
		log.Println("group already exists")
		return domain.Group{}, errors.New("group already exists")
	}
	groupRow, err := service.Create(ctx, group)
	if err != nil {
		log.Printf("failed to create group %v", err)
		return domain.Group{}, err
	}

	createdGroup, err := convertGroupsRowsToDomain(groupRow)
	if err != nil {
		log.Printf("failed to convert group into domain %v", err)
		return domain.Group{}, err
	}

	return createdGroup[0], err
}

func (repo *GroupServiceImpl) GetAll(ctx context.Context) ([]domain.Group, error) {
	service := repo.repo

	rows, err := service.GetAll(ctx)
	if err != nil {
		log.Printf("failed to get groups %v", err)
	}

	groups, err := convertGroupsRowsToDomain(rows)
	if err != nil {
		log.Printf("failed to convert groups into domain %v", err)

	}

	return groups, nil
}

func (repo *GroupServiceImpl) GetById(ctx context.Context, id int64) (domain.Group, error) {
	service := repo.repo

	row := service.GetById(ctx, id)

	group, err := convertGroupRowToDomain(row)
	if err != nil {
		log.Printf("failed to convert group into domain %v", err)

	}

	return group, nil
}

func (repo *GroupServiceImpl) Update(ctx context.Context, id int64, groupNumber string) (domain.Group, error) {
	service := repo.repo

	group := domain.Group{
		Id:          id,
		GroupNumber: groupNumber,
	}

	if !repo.isGroupExistsById(ctx, id) {
		log.Println("group doesn't exists")
		return domain.Group{}, errors.New("group doesn't exists")
	}

	groupRow, err := service.Update(ctx, group)
	if err != nil {
		log.Printf("failed to update group %v", err)
		return domain.Group{}, err
	}

	updatedGroup, err := convertGroupsRowsToDomain(groupRow)
	if err != nil {
		log.Printf("failed to convert group into domain %v", err)
		return domain.Group{}, err
	}

	return updatedGroup[0], err
}

func (repo *GroupServiceImpl) DeleteById(ctx context.Context, id int64) error {
	service := repo.repo

	if !repo.isGroupExistsById(ctx, id) {
		log.Println("group doesn't exists")
		return errors.New("group doesn't exist")
	}
	err := service.DeleteById(ctx, id)
	if err != nil {
		log.Printf("failed to delete group %v", err)
	}

	return nil
}

func convertGroupRowToDomain(row pgx.Row) (domain.Group, error) {
	var group domain.Group

	err := row.Scan(&group.Id, &group.GroupNumber)

	if err != nil {
		return domain.Group{}, err
	}

	return group, err
}

func convertGroupsRowsToDomain(rows pgx.Rows) ([]domain.Group, error) {
	var groups []domain.Group

	for rows.Next() {
		var r domain.Group
		err := rows.Scan(&r.Id, &r.GroupNumber)
		if err != nil {
			return nil, err
		}
		groups = append(groups, r)
	}

	return groups, nil
}

func (repo *GroupServiceImpl) isGroupExistsByNumber(ctx context.Context, groupNumber string) bool {
	service := repo.repo

	if errors.Is(service.GetByGroupNumber(ctx, groupNumber).Scan(), pgx.ErrNoRows) {
		return false
	}

	return true
}

func (repo *GroupServiceImpl) isGroupExistsById(ctx context.Context, id int64) bool {
	service := repo.repo

	if errors.Is(service.GetById(ctx, id).Scan(), pgx.ErrNoRows) {
		return false
	}

	return true
}
