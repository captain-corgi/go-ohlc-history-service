package entity

import "strconv"

// OHLC represent transaction info
//
//	UNIX,SYMBOL,OPEN,HIGH,LOW,CLOSE
type OHLC struct {
	Symbol string  `json:"symbol"      csv:"SYMBOL"  example:"BTCUSDT"`
	Id     int64   `json:"-"          csv:"-"       example:"1"`
	Unix   int64   `json:"unix"        csv:"UNIX"    example:"1644719700000"`
	Open   float64 `json:"open"        csv:"OPEN"    example:"42123.29000000"`
	High   float64 `json:"high"        csv:"HIGH"    example:"42148.32000000"`
	Low    float64 `json:"low"         csv:"LOW"     example:"42120.82000000"`
	Close  float64 `json:"close"       csv:"CLOSE"   example:"42146.06000000"`
}

func (r *OHLC) Columns() (cols []interface{}) {
	cols = append(cols, &r.Unix)
	cols = append(cols, &r.Symbol)
	cols = append(cols, &r.Open)
	cols = append(cols, &r.High)
	cols = append(cols, &r.Low)
	cols = append(cols, &r.Close)
	return
}

func (r *OHLC) Strings() []string {
	if r == nil {
		return []string{}
	}
	return []string{
		strconv.FormatInt(r.Unix, 10),
		r.Symbol,
		strconv.FormatFloat(r.Open, 'f', -1, 64),
		strconv.FormatFloat(r.High, 'f', -1, 64),
		strconv.FormatFloat(r.Low, 'f', -1, 64),
		strconv.FormatFloat(r.Close, 'f', -1, 64),
	}
}
