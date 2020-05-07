package amap

// DbPoi DbPoi
type DbPoi struct {
	Source string  `json:"source" xorm:"source"`
	PoiKey string  `json:"poikey" xorm:"poikey"`
	Base   Poi     `json:"base" xorm:"base"`
	Detail *Detail `json:"details" xorm:"details"`
	Common ` xorm:"extends"`
}

// TableName TableName
func (p DbPoi) TableName() string {
	return "poi"
}

// InsertPoi InsertPoi
func (h *Handler) InsertPoi(p DbPoi) (err error) {
	_, err = h.e.Insert(p)
	if err != nil {
		return err
	}
	return
}

// InsertPois InsertPois
func (h *Handler) InsertPois(p []DbPoi) (err error) {
	_, err = h.e.Insert(p)
	if err != nil {
		return err
	}
	return
}

// InsertPoi InsertPoi
func (h *Handler) findAllPoi() (pois []DbPoi, err error) {
	err = h.e.Find(&pois)
	return
}
