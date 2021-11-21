package database

import (
	"context"
	"time"

	"github.com/eurofurence/reg-auth-service/internal/repository/database/dbrepo"
	"github.com/eurofurence/reg-auth-service/internal/repository/database/inmemorydb"
	"github.com/eurofurence/reg-auth-service/internal/repository/logging"
)

var (
	ActiveRepository dbrepo.Repository
	pruneTicker      time.Ticker
	pruneStop        chan bool
)

// only exported so you can use it in test code - use Open()
func SetRepository(repository dbrepo.Repository) {
	ActiveRepository = repository
}

func Open() {
	var r dbrepo.Repository
	logging.NoCtx().Info("Opening inmemory database...")
	r = inmemorydb.Create()
	r.Open()
	pruneTicker = *time.NewTicker(3 * time.Second)
	go func() {
		for {
			select {
			case <-pruneStop:
				pruneTicker.Stop()
				return
			case <-pruneTicker.C:
				r.PruneAuthRequests(context.Background())
			}
		}
	}()
	SetRepository(r)
}

func Close() {
	logging.NoCtx().Info("Closing database...")
	pruneStop <- true
	GetRepository().Close()
	SetRepository(nil)
}

func GetRepository() dbrepo.Repository {
	if ActiveRepository == nil {
		logging.NoCtx().Fatal("You must Open() the database before using it. This is an error in your implementation.")
	}
	return ActiveRepository
}
