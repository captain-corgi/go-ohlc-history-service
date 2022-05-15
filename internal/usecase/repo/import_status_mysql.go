package repo

import (
	"context"
	"database/sql"
	"github.com/captain-corgi/go-ohlc-history-service/internal/entity"
	"github.com/captain-corgi/go-ohlc-history-service/pkg/logger"
)

//go:generate mockgen -source=import_status_mysql.go -destination=mocks/import_status_mysql_mock.go -package=mocks

type (
	ImportStatus interface {
		Save(ctx context.Context, ohlcModel entity.ImportStatus) (int64, error)
		FindByStatusAndCreatedDate(ctx context.Context, status string, createdDate string) (entity.ImportStatus, error)
	}
	ImportStatusImpl struct {
		logger logger.Logger
		db     *sql.DB
	}
)

func NewImportStatusRepo(db *sql.DB, logger logger.Logger) *ImportStatusImpl {
	return &ImportStatusImpl{db: db, logger: logger}
}

func (r *ImportStatusImpl) Save(ctx context.Context, ohlcModel entity.ImportStatus) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ImportStatusImpl) FindByStatusAndCreatedDate(ctx context.Context, status string, createdDate string) (entity.ImportStatus, error) {
	// TODO: [Anh Le] Implement
	return entity.ImportStatus{}, nil
}
