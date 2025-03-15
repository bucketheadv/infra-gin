package db

import (
	"database/sql"
	"github.com/bucketheadv/infra-core/modules/logger"
	"github.com/bucketheadv/infra-gin/api"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MySQLConf struct {
	Url string
}

func NewMySQL(config MySQLConf, gormConfig *gorm.Config) *gorm.DB {
	var err error
	DB, err := gorm.Open(mysql.Open(config.Url), gormConfig)
	if err != nil {
		logger.Fatal(err)
	}
	return DB
}

func Page[T schema.Tabler](tx *gorm.DB, page api.Page) (api.PageResult[T], error) {
	var data []T
	var model T
	var total int64
	tx.Model(model).Count(&total).Offset(page.Offset()).Limit(page.Limit()).Find(&data)

	var totalInt = (int)(total)
	var pages = 0
	if totalInt%page.PageSize == 0 {
		pages = totalInt / page.PageSize
	} else {
		pages = totalInt/page.PageSize + 1
	}
	return api.PageResult[T]{
		Page:    page,
		Pages:   pages,
		Total:   total,
		Records: data,
	}, nil
}

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		logger.Error(err)
	}
}
