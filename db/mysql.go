package db

import (
	"database/sql"
	"github.com/bucketheadv/infra-core/modules/logger"
	"github.com/bucketheadv/infra-gin"
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

func Page[T schema.Tabler](db *gorm.DB, page infra_gin.Page) (infra_gin.PageResult[T], error) {
	var tx = db.Offset(page.Offset()).Limit(page.Limit())
	var data []T
	tx.Find(&data)

	var total int64
	tx.Count(&total)
	var totalInt = (int)(total)
	var pages = 0
	if totalInt%page.PageSize == 0 {
		pages = totalInt / page.PageSize
	} else {
		pages = totalInt/page.PageSize + 1
	}
	return infra_gin.PageResult[T]{
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
