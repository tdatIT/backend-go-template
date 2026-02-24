package taskgroup

import (
	"context"

	"github.com/tdatIT/backend-go/internal/domain/models"
	"github.com/tdatIT/backend-go/pkgs/db/orm"
	"gorm.io/gorm"
)

type reposImpl struct {
	orm orm.ORM
}

func NewRepository(orm orm.ORM) Repository {
	return &reposImpl{
		orm: orm,
	}
}

func (r reposImpl) Create(ctx context.Context, item *models.TaskGroup) error {
	return r.orm.GormDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(item).Error
	})
}

func (r reposImpl) FindByID(ctx context.Context, id uint64) (*models.TaskGroup, error) {
	item := new(models.TaskGroup)
	err := r.orm.GormDB().
		WithContext(ctx).
		First(item, id).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r reposImpl) FindAllBy(ctx context.Context, params *GetListParams) ([]*models.TaskGroup, int64, error) {
	var (
		items []*models.TaskGroup
		count int64
	)

	if params == nil {
		params = &GetListParams{}
	}

	db := r.orm.GormDB().WithContext(ctx).Model(&models.TaskGroup{})

	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset(params.Offset).Limit(params.Limit).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, count, nil
}

func (r reposImpl) Update(ctx context.Context, item *models.TaskGroup) error {
	return r.orm.GormDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Save(item).Error
	})
}

func (r reposImpl) Delete(ctx context.Context, id uint64) error {
	return r.orm.GormDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.TaskGroup{}, id).Error
	})
}
