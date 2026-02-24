package task

import (
	"context"

	"github.com/tdatIT/backend-go/internal/domain/models"
)

type GetListParams struct {
	Offset int
	Limit  int
}

// Repository defines persistence operations for Task models.
type Repository interface {
	Create(ctx context.Context, item *models.Task) error
	FindByID(ctx context.Context, id uint64) (*models.Task, error)
	FindAllBy(ctx context.Context, params *GetListParams) ([]*models.Task, int64, error)
	Update(ctx context.Context, item *models.Task) error
	Delete(ctx context.Context, id uint64) error
}
