package taskgroup

import (
	"context"

	"github.com/tdatIT/backend-go/internal/domain/models"
)

type GetListParams struct {
	Offset int
	Limit  int
}

// Repository defines persistence operations for TaskGroup models.
type Repository interface {
	Create(ctx context.Context, item *models.TaskGroup) error
	FindByID(ctx context.Context, id uint64) (*models.TaskGroup, error)
	FindAllBy(ctx context.Context, params *GetListParams) ([]*models.TaskGroup, int64, error)
	Update(ctx context.Context, item *models.TaskGroup) error
	Delete(ctx context.Context, id uint64) error
}
