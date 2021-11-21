package inmemorydb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eurofurence/reg-auth-service/internal/entity"
	"github.com/eurofurence/reg-auth-service/internal/repository/database/dbrepo"
)

type InMemoryRepository struct {
	authRequests sync.Map
}

func Create() dbrepo.Repository {
	return &InMemoryRepository{}
}

func (r *InMemoryRepository) Open() {
	r.authRequests = sync.Map{}
}

func (r *InMemoryRepository) Close() {
	r.authRequests = sync.Map{}
}

func (r *InMemoryRepository) AddAuthRequest(ctx context.Context, ar *entity.AuthRequest) error {
	if _, ok := r.authRequests.Load(ar.State); ok {
		return fmt.Errorf("cannot add auth request '%s' - already present", ar.State)
	} else {
		// copy the entity, so later modifications won't also modify it in the in-memory db
		copiedEntity := *ar
		r.authRequests.Store(ar.State, &copiedEntity)
		return nil
	}
}

func (r *InMemoryRepository) GetAuthRequestByState(ctx context.Context, state string) (*entity.AuthRequest, error) {
	if ar, ok := r.authRequests.Load(state); ok {
		if ar.(*entity.AuthRequest).ExpiresAt.Before(time.Now()) {
			r.authRequests.Delete(state)
			return nil, fmt.Errorf("cannot get auth request '%s' - already expired", state)
		} else {
			// copy the entity, so later modifications won't also modify it in the in-memory db
			copiedEntity := *ar.(*entity.AuthRequest)
			return &copiedEntity, nil
		}
	} else {
		return nil, fmt.Errorf("cannot get auth request '%s' - not present", state)
	}
}

func (r *InMemoryRepository) DeleteAuthRequestByState(ctx context.Context, state string) error {
	if _, ok := r.authRequests.LoadAndDelete(state); ok {
		return nil
	} else {
		return fmt.Errorf("cannot delete auth request '%s' - not present", state)
	}
}

func (r *InMemoryRepository) PruneAuthRequests(ctx context.Context) (uint, error) {
	pruneCount := uint(0)

	r.authRequests.Range(func(state, ar interface{}) bool {
		if ar.(*entity.AuthRequest).ExpiresAt.Before(time.Now()) {
			r.authRequests.Delete(state)
			pruneCount++
		}
		return true
	})

	return pruneCount, nil
}
