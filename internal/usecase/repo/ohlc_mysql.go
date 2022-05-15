package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/captain-corgi/go-ohlc-history-service/internal/entity"
	"github.com/captain-corgi/go-ohlc-history-service/pkg/logger"
	"strings"
)

//go:generate mockgen -source=ohlc_mysql.go -destination=mocks/ohlc_mysql_mock.go -package=mocks

type (
	OHLC interface {
		SaveAll(ctx context.Context, ohlcModels []entity.OHLC) (int64, error)
		Count(ctx context.Context, searchModel entity.OHLCSearchRequest) (int64, error)
		Search(ctx context.Context, searchModel entity.OHLCSearchRequest) ([]entity.OHLC, error)
	}
	OHLCImpl struct {
		logger logger.Logger
		db     *sql.DB
	}
)

func NewOHLCRepo(db *sql.DB, logger logger.Logger) *OHLCImpl {
	return &OHLCImpl{db: db, logger: logger}
}

func (r *OHLCImpl) SaveAll(ctx context.Context, ohlcModels []entity.OHLC) (int64, error) {
	var (
		valueStrings []string
		valueArgs    []interface{}
	)
	for _, ohlcModel := range ohlcModels {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, ohlcModel.Columns()...)
	}

	smt := `INSERT IGNORE INTO ohlc (unix, symbol, open, high, low, close) VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	// r.logger.Debug("query: %s\n", smt) // NOTE: Turn on for more debugging info

	result, err := r.db.ExecContext(ctx, smt, valueArgs...)
	if err != nil {
		return 0, fmt.Errorf("save: %v", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("save: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (r *OHLCImpl) Count(ctx context.Context, searchModel entity.OHLCSearchRequest) (int64, error) {
	// Build conditions
	valueArgs, conditions := searchModel.BuildSearchCondition()

	// Query database
	query := fmt.Sprintf(`SELECT COUNT(*) FROM ohlc WHERE %s`, conditions)
	r.logger.Debug("query: %s\n", query)
	rows, err := r.db.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		return 0, fmt.Errorf("count: %v", err)
	}

	var count int64
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, fmt.Errorf("count: %v", err)
		}
	}
	if rows.Err() != nil {
		return 0, fmt.Errorf("count: %v", err)
	}
	return count, nil
}

func (r *OHLCImpl) Search(ctx context.Context, searchModel entity.OHLCSearchRequest) ([]entity.OHLC, error) {
	// Build conditions
	valueArgs, conditions := searchModel.BuildSearchCondition()
	limitOffsetArgs, limitOffset := searchModel.BuildLimitOffset()
	valueArgs = append(valueArgs, limitOffsetArgs...)

	// Query database
	query := fmt.Sprintf(`SELECT unix, symbol, open, high, low, close FROM ohlc WHERE %s %s`, conditions, limitOffset)
	r.logger.Debug("query: %s\n", query)
	rows, err := r.db.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		return []entity.OHLC{}, fmt.Errorf("select: %v", err)
	}
	defer rows.Close()

	var results []entity.OHLC
	for rows.Next() {
		var result entity.OHLC
		if err = rows.Scan(result.Columns()...); err != nil {
			return []entity.OHLC{}, fmt.Errorf("select: %v", err)
		}
		results = append(results, result)
	}
	if rows.Err() != nil {
		return []entity.OHLC{}, fmt.Errorf("select: %v", err)
	}
	return results, nil
}
