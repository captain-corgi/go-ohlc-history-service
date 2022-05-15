package entity

type ImportStatus struct {
	Id        int64  `json:"id"`
	OHLCId    string `json:"OHLCId"`
	Status    string `json:"status"`
	Count     int64  `json:"count"`
	Reason    string `json:"reason"`
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
