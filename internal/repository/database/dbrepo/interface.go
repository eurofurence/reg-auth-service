package dbrepo

import (
	"context"

	"github.com/eurofurence/reg-auth-service/internal/entity"
)

type Repository interface {
	Open() error
	Close()

	AddAuthRequest(ctx context.Context, ar *entity.AuthRequest) error
	GetAuthRequestByState(ctx context.Context, state string) (*entity.AuthRequest, error)
	DeleteAuthRequestByState(ctx context.Context, state string) error

	PruneAuthRequests(ctx context.Context) (uint, error)
}
