package user

import (
	"context"

	"github.com/tdatIT/backend-go/internal/domain/models"
)

type GetListParams struct {
	Offset int
	Limit  int
}

// Repository defines persistence operations for User models.
type Repository interface {
	Create(ctx context.Context, item *models.User) error
	FindByID(ctx context.Context, id uint64) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByOIDC(ctx context.Context, provider string, subject string) (*models.User, error)
	FindAllAndCount(ctx context.Context, params GetListParams) ([]*models.User, int64, error)
	Update(ctx context.Context, item *models.User) error
	Delete(ctx context.Context, id uint64) error
}
