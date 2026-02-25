package user

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

func (r reposImpl) Create(ctx context.Context, item *models.User) error {
	return r.orm.GormDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(item).Error
	})
}

func (r reposImpl) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	item := new(models.User)
	err := r.orm.GormDB().
		WithContext(ctx).
		First(item, id).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r reposImpl) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	item := new(models.User)
	err := r.orm.GormDB().
		WithContext(ctx).
		Where("username = ?", username).
		First(item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r reposImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	item := new(models.User)
	err := r.orm.GormDB().
		WithContext(ctx).
		Where("email = ?", email).
		First(item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r reposImpl) FindByOIDC(ctx context.Context, provider string, subject string) (*models.User, error) {
	item := new(models.User)
	err := r.orm.GormDB().
		WithContext(ctx).
		Where("oidc_provider = ? AND oidc_subject = ?", provider, subject).
		First(item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r reposImpl) FindAllAndCount(ctx context.Context, params GetListParams) ([]*models.User, int64, error) {
	var (
		items []*models.User
		count int64
	)

	db := r.orm.GormDB().WithContext(ctx).Model(&models.User{})

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

func (r reposImpl) Update(ctx context.Context, item *models.User) error {
	return r.orm.GormDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Save(item).Error
	})
}

func (r reposImpl) Delete(ctx context.Context, id uint64) error {
	return r.orm.GormDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.User{}, id).Error
	})
}
