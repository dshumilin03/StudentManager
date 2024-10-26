package service

import (
	"StudentManager/internal/custom_errors"
	"StudentManager/internal/domain"
	"StudentManager/internal/dto"
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
	ctx context.Context, groupDto dto.GroupDto) (domain.Group, error) {
	service := repo.repo

	group := domain.Group{
		GroupNumber: groupDto.GroupNumber,
	}

	if repo.IsGroupExistsByNumber(ctx, group.GroupNumber) {

		log.Println("group already exists")
		return domain.Group{}, custom_errors.ErrGroupExists
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

	log.Printf("created group: %v", createdGroup)
	return createdGroup[0], err
}

func (repo *GroupServiceImpl) GetAll(ctx context.Context) ([]domain.Group, error) {
	service := repo.repo

	rows, err := service.GetAll(ctx)
	if err != nil {
		log.Printf("failed to get groups %v", err)
		return []domain.Group{}, err
	}

	groups, err := convertGroupsRowsToDomain(rows)
	if err != nil {
		log.Printf("failed to convert groups into domain %v", err)

		return []domain.Group{}, err
	}

	log.Println("received all groups")

	return groups, nil
}

func (repo *GroupServiceImpl) GetById(ctx context.Context, id int64) (domain.Group, error) {
	service := repo.repo

	row := service.GetById(ctx, id)

	group, err := convertGroupRowToDomain(row)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("group does not exist")
			return domain.Group{}, custom_errors.ErrGroupNotFound
		}

		log.Printf("failed to convert group into domain %v", err)
		return domain.Group{}, err
	}

	log.Printf("received group by id: %v", group)

	return group, nil
}

func (repo *GroupServiceImpl) Update(ctx context.Context,
	groupDto dto.GroupDto) (domain.Group, error) {
	service := repo.repo

	group := domain.Group{
		Id:          groupDto.Id,
		GroupNumber: groupDto.GroupNumber,
	}

	if !repo.IsGroupExistsById(ctx, group.Id) {
		log.Println("group does not exist")
		return domain.Group{}, custom_errors.ErrGroupNotFound
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

	log.Printf("group updated: %v", updatedGroup)
	return updatedGroup[0], err
}

func (repo *GroupServiceImpl) DeleteById(ctx context.Context, id int64) error {
	service := repo.repo

	if !repo.IsGroupExistsById(ctx, id) {
		log.Println("group does not exist")
		return custom_errors.ErrGroupNotFound
	}
	err := service.DeleteById(ctx, id)
	if err != nil {
		log.Printf("failed to delete group %v", err)

		return err
	}

	log.Printf("deleted group with id: %v", id)
	return nil
}

func convertGroupRowToDomain(row pgx.Row) (domain.Group, error) {
	var group domain.Group

	err := row.Scan(&group.Id, &group.GroupNumber)

	if err != nil {
		return domain.Group{}, err
	}

	log.Println("successfully converted group row to domain")

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

	log.Println("successfully converted groups rows to domain")

	return groups, nil
}

func (repo *GroupServiceImpl) IsGroupExistsByNumber(ctx context.Context, groupNumber string) bool {
	service := repo.repo

	err := service.GetByGroupNumber(ctx, groupNumber)

	if errors.Is(err.Scan(), pgx.ErrNoRows) {
		return false
	}

	return true
}

func (repo *GroupServiceImpl) IsGroupExistsById(ctx context.Context, id int64) bool {
	service := repo.repo

	if errors.Is(service.GetById(ctx, id).Scan(), pgx.ErrNoRows) {
		return false
	}

	return true
}

func (g *GroupServiceImpl) GetService() Service[dto.GroupDto, domain.Group] {
	return g
}
