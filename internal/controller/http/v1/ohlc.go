package v1

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/captain-corgi/go-ohlc-history-service/config"
	"github.com/captain-corgi/go-ohlc-history-service/pkg/cache"
	"github.com/labstack/echo/v4"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/captain-corgi/go-ohlc-history-service/internal/entity"
	"github.com/captain-corgi/go-ohlc-history-service/internal/usecase"
	"github.com/captain-corgi/go-ohlc-history-service/pkg/logger"
)

// Represent DB state
var (
	InMemoryDatabase = make([]entity.OHLC, 0)
	// WriteThreshold is the number of rows for commit data and clear memory
	WriteThreshold = 10000
	// ProcessTimeout is the timeout for processing request
	ProcessTimeout = time.Duration(15) * time.Minute
)

type ohlcRoutes struct {
	useCase usecase.OHLCUseCase
	logger  logger.Interface
}

func newOHLCRoutes(handler *echo.Group, useCase usecase.OHLCUseCase, logger logger.Interface, cfg *config.Config) {
	// Init cache
	cache.InMemory.Set("status", entity.ProcessStatus_NEW)
	cache.InMemory.Set("message", "")
	cache.InMemory.Set("failedRecords", 0)

	// Init default values
	if cfg.MySQL.WriteThreshold != 0 {
		WriteThreshold = cfg.MySQL.WriteThreshold
	}
	if cfg.HTTP.ProcessTimeout != 0 {
		ProcessTimeout = time.Duration(cfg.HTTP.ProcessTimeout) * time.Minute
	}

	// Init routers
	r := &ohlcRoutes{useCase, logger}
	h := handler.Group("/ohlc")
	{
		h.GET("/data", r.getData)
		h.POST("/data", r.postData)
	}
}

// @Summary     Search OHLC data
// @Description Show OHLC data matching the given criteria
// @ID          get_ohlc
// @Tags  	    ohlc
// @Accept      json
// @Produce     json
// @Param       unix         query int    false "1644719700000"    "unix timestamp"
// @Param       symbol       query string false "BTCUSDT"          "symbol"
// @Param       open         query number false "42123.29000000"   "open"
// @Param       high         query number false "42148.32000000"   "high"
// @Param       low          query number false "42120.82000000"   "low"
// @Param       close        query number false "42146.06000000"   "close"
// @Param       page         query int    false "1"                "page"
// @Param       itemsPerPage query int    false "20"               "items per page"
// @Success     200 {object} entity.OHLCSearchResponse
// @Failure     500 {object} response
// @Router      /ohlc/data [get]
func (r *ohlcRoutes) getData(c echo.Context) error {
	// Prepare
	var (
		err error
		req entity.OHLCSearchRequest
		res entity.OHLCSearchResponse
	)
	Status, ok := cache.InMemory.Get("status")
	if !ok {
		Status = entity.ProcessStatus_NEW
	}
	Message, ok := cache.InMemory.Get("message")
	if !ok {
		Message = ""
	}
	var TimerInt int64
	Timer, ok := cache.InMemory.Get("timer")
	if ok {
		TimerInt, _ = strconv.ParseInt(fmt.Sprintf("%d", Timer), 10, 64)
	} else {
		TimerInt = 0
	}
	FailedRecords, ok := cache.InMemory.Get("failedRecords")
	if !ok {
		FailedRecords = 0
	}

	if Status == entity.ProcessStatus_FAIL && Message != "" {
		return errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("process failed: %s", Message))
	}

	// Bind search condition
	formParams, _ := c.FormParams()
	req.Unix, _ = strconv.ParseInt(formParams.Get("unix"), 10, 64)
	req.Symbol = formParams.Get("symbol")
	req.Open, _ = strconv.ParseFloat(formParams.Get("open"), 64)
	req.High, _ = strconv.ParseFloat(formParams.Get("high"), 64)
	req.Low, _ = strconv.ParseFloat(formParams.Get("low"), 64)
	req.Close, _ = strconv.ParseFloat(formParams.Get("close"), 64)
	req.Page, _ = strconv.Atoi(formParams.Get("page"))
	req.ItemsPerPage, _ = strconv.Atoi(formParams.Get("itemsPerPage"))

	// Search
	res, err = r.useCase.Search(c.Request().Context(), req)
	if err != nil {
		r.logger.Error(err, "http - v1 - history")
		return errorResponse(c, http.StatusInternalServerError, "database problems")
	}
	res.LatestProcess.Status = entity.ProcessStatus(fmt.Sprintf("%s", Status))
	res.LatestProcess.Message = fmt.Sprintf("%s", Message)
	res.LatestProcess.LastProcessTime = fmt.Sprintf("%dms", TimerInt)
	FailedRecordsInt, _ := strconv.ParseInt(fmt.Sprintf("%d", FailedRecords), 10, 64)
	res.LatestProcess.FailedRecords = FailedRecordsInt

	// Response
	return c.JSON(http.StatusOK, res)
}

// @Summary     Upload OHLC data
// @Description Upload OHLC data to database
// @ID          post_ohlc
// @Tags  	    ohlc
// @Accept      json
// @Produce     json
// @Param       csv formData file true "csv file"
// @Success     200 {object} entity.OHLCSaveResponse
// @Failure     500 {object} response
// @Router      /ohlc/data [post]
func (r *ohlcRoutes) postData(c echo.Context) error {
	// Prepare
	var (
		err error
		res entity.OHLCSaveResponse
	)
	Status, ok := cache.InMemory.Get("status")
	if !ok {
		Status = entity.ProcessStatus_NEW
	}

	// If there are any process being processing, stop application from receive new file
	if Status == entity.ProcessStatus_PROCESSING {
		return errorResponse(c, http.StatusBadRequest, "Others file is being processing...")
	}

	// Bind search condition
	csvFile, err := c.FormFile("csv")
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, err.Error())
	}
	if !strings.HasSuffix(csvFile.Filename, ".csv") {
		return errorResponse(c, http.StatusBadRequest, "Only csv file is allowed")
	}

	// Process
	Status = entity.ProcessStatus_PROCESSING
	res.Status = entity.ProcessStatus_PROCESSING
	cache.InMemory.Set("status", Status)
	cache.InMemory.Set("timer", 0)
	go Process(r, csvFile)

	// Response
	return c.JSON(http.StatusOK, res)
}

func Process(r *ohlcRoutes, csvFile *multipart.FileHeader) {
	start := time.Now()
	r.logger.Info("Processing...")
	defer r.logger.Info("Finished")

	ctx, cancel := context.WithTimeout(context.Background(), ProcessTimeout)
	defer cancel()

	FailedRecords := 0

	// Open csv source
	src, err := csvFile.Open()
	if err != nil {
		cache.InMemory.Set("status", entity.ProcessStatus_FAIL)
		cache.InMemory.Set("message", fmt.Sprintf("failed to open csv file: %s", err.Error()))
		WriteHistory(ctx, r, 0, fmt.Sprintf("failed to open csv file: %s", err.Error()))
		return
	}
	defer src.Close()

	// Write output
	scanner := bufio.NewScanner(src)
	// optionally, resize scanner's capacity for lines over 64K, do it later
	for scanner.Scan() {
		reader := csv.NewReader(bytes.NewReader(scanner.Bytes()))
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				cache.InMemory.Set("status", entity.ProcessStatus_FAIL)
				cache.InMemory.Set("message", fmt.Sprintf("Bad data in this line: %s. Skipped.", err.Error()))
				FailedRecords++
				r.logger.Error("Bad data in this line: %s. Skipped.", err.Error())
				continue // NOTE: Continue to process next line
			}

			// Validate CSV file
			if len(record) != 6 {
				// Bad data will result empty line.
				FailedRecords++
				record = []string{"", "", "", "", "", ""}
			}
			var (
				unix                 int64
				symbol               string
				open, high, low, cls float64
				errs                 = make([]error, 0)
			)
			if unix, err = strconv.ParseInt(record[0], 10, 64); err != nil {
				errs = append(errs, err)
			}
			symbol = record[1]
			if open, err = strconv.ParseFloat(record[2], 64); err != nil {
				errs = append(errs, err)
			}
			if high, err = strconv.ParseFloat(record[3], 64); err != nil {
				errs = append(errs, err)
			}
			if low, err = strconv.ParseFloat(record[4], 64); err != nil {
				errs = append(errs, err)
			}
			if cls, err = strconv.ParseFloat(record[5], 64); err != nil {
				errs = append(errs, err)
			}
			if len(errs) > 0 {
				cache.InMemory.Set("status", entity.ProcessStatus_FAIL)
				cache.InMemory.Set("message", fmt.Sprintf("Bad data in this line. Skipped."))
				FailedRecords++
				r.logger.Error("Bad data in this line. Skipped.")
				continue // NOTE: Continue to process next line
			}

			InMemoryDatabase = append(InMemoryDatabase, entity.OHLC{
				Unix:   unix,
				Symbol: symbol,
				Open:   open,
				High:   high,
				Low:    low,
				Close:  cls,
			})

			// Represent store DB and clean up memory
			if len(InMemoryDatabase) > WriteThreshold {
				failedRecords, message := PersistAndClean(ctx, r)
				WriteHistory(ctx, r, failedRecords, message)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		cache.InMemory.Set("status", entity.ProcessStatus_FAIL)
		cache.InMemory.Set("message", fmt.Sprintf("failed to scan csv file: %s", err.Error()))
		return
	}

	// Store failed records
	cache.InMemory.Set("failedRecords", FailedRecords)

	PersistAndClean(ctx, r)
	cache.InMemory.Set("status", entity.ProcessStatus_SUCCESS)
	cache.InMemory.Set("message", "")

	end := time.Since(start).Milliseconds()
	cache.InMemory.Set("timer", end)
	r.logger.Info("Process took %dms", end)
}

func WriteHistory(ctx context.Context, r *ohlcRoutes, records int, message string) {
	r.logger.Info("Writing history...")
	defer r.logger.Info("Finished")

	// TODO: [Anh Le] Implement write history
	return
}

func PersistAndClean(ctx context.Context, r *ohlcRoutes) (int, string) {
	Data := InMemoryDatabase
	InMemoryDatabase = make([]entity.OHLC, 0)
	defer func() {
		Data = nil
	}()

	// Persist data to DB
	affectedRows, err := r.useCase.Save(ctx, Data)
	if err != nil {
		r.logger.Error(err.Error())
		return len(Data), err.Error()
	}

	if affectedRows != int64(len(Data)) {
		return int(int64(len(Data)) - affectedRows), "Failed to persist data"
	}

	return 0, ""
}
