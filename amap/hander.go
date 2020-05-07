package amap

import (
	"xorm.io/xorm"
)

// Handler Handler
type Handler struct {
	e *xorm.Engine
}

// SyncDb SyncDb
func (h *Handler) SyncDb() (err error) {
	p := new(DbPoi)

	err = h.e.Sync2(p)

	if err != nil {
		return err
	}

	return
}

// NewHandler NewHandler
func NewHandler(e *xorm.Engine) Handler {
	return Handler{e}
}

// Common Common
type Common struct {
	ID        int64 `json:"id" xorm:"id pk notnull  autoincr"`
	UpdatedAt Time  `json:"updated" xorm:"updated comment('创建时间')" `
	CreatedAt Time  `json:"created" xorm:"created comment('更新时间')"`
}
