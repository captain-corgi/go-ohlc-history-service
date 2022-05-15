package entity

// OHLCSearchRequest is the request for search OHLC data.
type OHLCSearchRequest struct {
	OHLC
	Page         int `json:"page"`
	ItemsPerPage int `json:"itemsPerPage"`
}

// BuildSearchCondition Return search condition based on input value
func (searchModel OHLCSearchRequest) BuildSearchCondition() (valueArgs []interface{}, conditions string) {
	valueArgs = make([]interface{}, 0)
	conditions = ` 1 = 1 `

	if searchModel.Unix != 0 {
		conditions += ` AND ohlc.unix = ? `
		valueArgs = append(valueArgs, searchModel.Unix)
	}
	if searchModel.Symbol != "" {
		conditions += ` AND ohlc.symbol = ? `
		valueArgs = append(valueArgs, searchModel.Symbol)
	}
	if searchModel.Open != 0 {
		conditions += ` AND ohlc.open = ? `
		valueArgs = append(valueArgs, searchModel.Open)
	}
	if searchModel.High != 0 {
		conditions += ` AND ohlc.open = ? `
		valueArgs = append(valueArgs, searchModel.High)
	}
	if searchModel.Low != 0 {
		conditions += ` AND ohlc.open = ? `
		valueArgs = append(valueArgs, searchModel.Low)
	}
	if searchModel.Close != 0 {
		conditions += ` AND ohlc.open = ? `
		valueArgs = append(valueArgs, searchModel.Close)
	}

	return
}

func (searchModel OHLCSearchRequest) BuildLimitOffset() (valueArgs []interface{}, sql string) {
	valueArgs = make([]interface{}, 0)
	sql = " LIMIT ? OFFSET ? "

	if searchModel.ItemsPerPage > 0 {
		valueArgs = append(valueArgs, searchModel.ItemsPerPage)
	} else {
		valueArgs = append(valueArgs, 10)
	}

	if searchModel.Page > 0 {
		valueArgs = append(valueArgs, (searchModel.Page-1)*searchModel.ItemsPerPage)
	} else {
		valueArgs = append(valueArgs, 0)
	}

	return
}

// OHLCSearchResponse is the response for search OHLC data.
type OHLCSearchResponse struct {
	LatestProcess ProcessingStatus `json:"latestProcess"`
	Offset        int              `json:"offset"`
	Limit         int              `json:"limit"`
	Total         int64            `json:"total"`
	Data          []OHLC           `json:"data"`
}

type ProcessingStatus struct {
	Status          ProcessStatus `json:"status"`
	Message         string        `json:"message"`
	FailedRecords   int64         `json:"failedRecords"`
	LastProcessTime string        `json:"lastProcessTime"`
}

type OHLCSaveResponse struct {
	Status ProcessStatus `json:"status"`
}
