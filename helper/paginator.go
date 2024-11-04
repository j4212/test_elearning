package helper

import (
	"math"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Pagination struct {
	TotalPages  int   `json:"total_pages"`
	Limit       int   `json:"limit"`
	CurrentPage int   `json:"current_page"`
	TotalRows   int64 `json:"total_rows"`
}

func Paginator(page int, limit int, tx *gorm.DB, rowFromTx int64) (*Pagination, *gorm.DB) {
	logrus.Infof("[repository] Passed arguments Page:%d and Limit:%d \n", page, limit)

	logrus.Infoln("[helper-paginator] Row From Tx:", rowFromTx)

	totalPages := int(math.Ceil(float64(rowFromTx) / float64(limit)))

	pageInformation := Pagination{
		TotalPages:  totalPages,
		TotalRows:   rowFromTx,
		CurrentPage: page,
		Limit:       limit,
	}

	p := tx.Scopes(func(d *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return d.Offset(offset).Limit(limit)
	})

	return &pageInformation, p
}
