package database

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"time"

	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/database/dbrepo"
	"github.com/eurofurence/reg-auth-service/internal/repository/database/inmemorydb"
)

var (
	ActiveRepository dbrepo.Repository
	pruneTicker      *time.Ticker
	pruneStop        chan bool
)

// only exported so you can use it in test code - use Open()
func SetRepository(repository dbrepo.Repository) {
	ActiveRepository = repository
}

func Open() error {
	var r dbrepo.Repository
	aulogging.Logger.NoCtx().Info().Print("Opening inmemory database...")
	r = inmemorydb.Create()
	r.Open()
	pruneTicker = time.NewTicker(config.AuthRequestTimeout())
	pruneStop = make(chan bool)
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
	return nil
}

func Close() {
	aulogging.Logger.NoCtx().Info().Print("Closing database...")
	pruneStop <- true
	GetRepository().Close()
	SetRepository(nil)
}

func GetRepository() dbrepo.Repository {
	if ActiveRepository == nil {
		aulogging.Logger.NoCtx().Fatal().Print("You must Open() the database before using it. This is an error in your implementation.")
	}
	return ActiveRepository
}
