package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/captain-corgi/go-ohlc-history-service/internal/entity"
	"github.com/captain-corgi/go-ohlc-history-service/internal/usecase/repo"
	"github.com/captain-corgi/go-ohlc-history-service/pkg/logger"
	"time"
)

// In-memory flag
var (
	InMemoryDatabase = make([]entity.OHLC, 0)
	Count            = 0
	Message          = ""
	Status           = entity.ProcessStatus_NEW
	FailedRecords    int
)

// OHLCUseCase - play the role of service
type (
	OHLC interface {
		Save(ctx context.Context, ohlcList []entity.OHLC) (affected int64, err error)
		Search(ctx context.Context) ([]entity.OHLC, error)
	}
	OHLCUseCase struct {
		logger           logger.Logger
		ohlcRepo         repo.OHLC
		importStatusRepo repo.ImportStatus
	}
)

// NewOHLCUseCase -.
func NewOHLCUseCase(db *sql.DB, logger logger.Logger) *OHLCUseCase {
	return &OHLCUseCase{
		logger:           logger,
		ohlcRepo:         repo.NewOHLCRepo(db, logger),
		importStatusRepo: repo.NewImportStatusRepo(db, logger),
	}
}

// Save receives OHLC data from web API.
func (r *OHLCUseCase) Save(ctx context.Context, ohlcList []entity.OHLC) (affected int64, err error) {
	// Greeting
	r.logger.Info(fmt.Sprintf("Saving a bulk of %d rows OHLC...", len(ohlcList)))

	// Save OHLC data
	affected, err = r.ohlcRepo.SaveAll(ctx, ohlcList)
	if err != nil {
		return 0, err
	}

	// Return
	r.logger.Info(fmt.Sprintf("Saved %d rows OHLC data.", affected))
	return affected, nil
}

// Search searches OHLC data from repository.
func (r *OHLCUseCase) Search(ctx context.Context, searchModel entity.OHLCSearchRequest) (entity.OHLCSearchResponse, error) {
	// Greeting
	r.logger.Info("Searching OHLC data...")
	defer r.logger.Info("Search OHLC data completed.")

	// Count OHLC data
	count, err := r.ohlcRepo.Count(ctx, searchModel)
	if err != nil {
		return entity.OHLCSearchResponse{}, err
	}

	// Search OHLC data
	ohlcList, err := r.ohlcRepo.Search(ctx, searchModel)
	if err != nil {
		return entity.OHLCSearchResponse{}, err
	}

	// Search failed records
	today := time.Now().Format(entity.DateFormat)
	failedRecords, err := r.importStatusRepo.FindByStatusAndCreatedDate(ctx, string(entity.ProcessStatus_FAIL), today)
	if err != nil {
		return entity.OHLCSearchResponse{}, err
	}

	// Return
	return entity.OHLCSearchResponse{
		LatestProcess: entity.ProcessingStatus{
			Status:        entity.ProcessStatus_SUCCESS,
			Message:       failedRecords.Reason,
			FailedRecords: failedRecords.Count,
		},
		Offset: searchModel.Page,
		Limit:  searchModel.ItemsPerPage,
		Total:  count,
		Data:   ohlcList,
	}, nil
}
