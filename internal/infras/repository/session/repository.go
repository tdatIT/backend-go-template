package session

import (
	"context"

	"github.com/tdatIT/backend-go/internal/domain/models"
)

// Repository defines persistence operations for Session models.
type Repository interface {
	Create(ctx context.Context, item *models.Session) error
	FindByID(ctx context.Context, id string) (*models.Session, error)
	FindByRefreshJTI(ctx context.Context, jti string) (*models.Session, error)
	FindBySessionID(ctx context.Context, sessionID string) (*models.Session, error)
	Update(ctx context.Context, item *models.Session) error
	RotateRefreshJTI(ctx context.Context, id string, oldJTI string, newJTI string) error
	Deactivate(ctx context.Context, id string) error
}
