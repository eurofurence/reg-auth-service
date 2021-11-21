package database

import (
	"github.com/eurofurence/reg-auth-service/internal/repository/database/dbrepo"
	"github.com/eurofurence/reg-auth-service/internal/repository/database/inmemorydb"
	"github.com/eurofurence/reg-auth-service/internal/repository/logging"
)

var (
	ActiveRepository dbrepo.Repository
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
	SetRepository(r)
}

func Close() {
	logging.NoCtx().Info("Closing database...")
	GetRepository().Close()
	SetRepository(nil)
}

func GetRepository() dbrepo.Repository {
	if ActiveRepository == nil {
		logging.NoCtx().Fatal("You must Open() the database before using it. This is an error in your implementation.")
	}
	return ActiveRepository
}
