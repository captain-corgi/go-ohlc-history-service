package entity

type ImportStatus struct {
	OHLCId    string `json:"OHLCId"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
	Id        int64  `json:"id"`
	Count     int64  `json:"count"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (r *ImportStatus) Columns() (cols []interface{}) {
	cols = append(cols, &r.Id)
	cols = append(cols, &r.OHLCId)
	cols = append(cols, &r.Status)
	cols = append(cols, &r.Count)
	cols = append(cols, &r.Reason)
	cols = append(cols, &r.CreatedAt)
	cols = append(cols, &r.UpdatedAt)
	return
}
