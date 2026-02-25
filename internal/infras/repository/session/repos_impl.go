package session

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/tdatIT/backend-go/internal/domain/models"
	"github.com/tdatIT/backend-go/pkgs/cache"
	"github.com/tdatIT/backend-go/pkgs/db/orm"
	"github.com/tdatIT/backend-go/pkgs/utils/genid"
	"gorm.io/gorm"
)

type reposImpl struct {
	orm        orm.ORM
	cache      cache.Cache
	sessionTTL time.Duration
}

func NewRepository(orm orm.ORM, cacheClient cache.Cache, sessionTTL time.Duration) Repository {
	return &reposImpl{
		orm:        orm,
		cache:      cacheClient,
		sessionTTL: sessionTTL,
	}
}

func (r reposImpl) Create(ctx context.Context, item *models.Session) error {
	if item.ID == "" {
		item.ID = genid.GenerateNanoID()
	}

	err := r.orm.GormDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(item).Error
	})
	if err != nil {
		return err
	}

	return r.setCache(ctx, item)
}

func (r reposImpl) FindByID(ctx context.Context, id string) (*models.Session, error) {
	if cached, err := r.getCache(ctx, r.cacheKeyByID(id)); err == nil {
		return cached, nil
	}

	item := new(models.Session)
	err := r.orm.GormDB().WithContext(ctx).First(item, id).Error
	if err != nil {
		return nil, err
	}

	_ = r.setCache(ctx, item)
	return item, nil
}

func (r reposImpl) FindByRefreshJTI(ctx context.Context, jti string) (*models.Session, error) {
	item := new(models.Session)
	err := r.orm.GormDB().WithContext(ctx).
		Where("refresh_jti = ?", jti).
		First(item).Error
	if err != nil {
		return nil, err
	}

	_ = r.setCache(ctx, item)
	return item, nil
}

func (r reposImpl) FindBySessionID(ctx context.Context, sessionID string) (*models.Session, error) {
	if cached, err := r.getCache(ctx, r.cacheKeyByID(sessionID)); err == nil {
		return cached, nil
	}

	item := new(models.Session)
	err := r.orm.GormDB().WithContext(ctx).First(item, sessionID).Error
	if err != nil {
		return nil, err
	}

	_ = r.setCache(ctx, item)
	return item, nil
}

func (r reposImpl) Update(ctx context.Context, item *models.Session) error {
	err := r.orm.GormDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Save(item).Error
	})
	if err != nil {
		return err
	}

	return r.setCache(ctx, item)
}

func (r reposImpl) RotateRefreshJTI(ctx context.Context, id string, oldJTI string, newJTI string) error {
	item := new(models.Session)
	if err := r.orm.GormDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(item).Error; err != nil {
			return err
		}

		if item.RefreshJTI != oldJTI {
			return gorm.ErrRecordNotFound
		}

		item.RefreshJTI = newJTI
		item.LastUsedAt = new(time.Now())
		return tx.Save(item).Error
	}); err != nil {
		return err
	}

	return r.setCache(ctx, item)
}

func (r reposImpl) Deactivate(ctx context.Context, id string) error {
	item := new(models.Session)
	if err := r.orm.GormDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(item, id).Error; err != nil {
			return err
		}
		item.IsActive = false
		now := time.Now()
		item.LastUsedAt = &now
		return tx.Save(item).Error
	}); err != nil {
		return err
	}

	return r.setCache(ctx, item)
}

func (r reposImpl) cacheKeyByID(id string) string {
	return fmt.Sprintf("session:id:%s", id)
}

func (r reposImpl) getCache(ctx context.Context, key string) (*models.Session, error) {
	if r.cache == nil {
		return nil, fmt.Errorf("cache disabled")
	}

	data, err := r.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	item := new(models.Session)
	if err := sonic.Unmarshal(data, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (r reposImpl) setCache(ctx context.Context, item *models.Session) error {
	data, err := sonic.Marshal(item)
	if err != nil {
		return err
	}

	return r.cache.Set(ctx, r.cacheKeyByID(item.ID), data, r.sessionTTL)
}
