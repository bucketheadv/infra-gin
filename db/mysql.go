package db

import (
	"database/sql"
	"github.com/bucketheadv/infra-core/modules/logger"
	"github.com/bucketheadv/infra-gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func Page(db *gorm.DB, page infra_gin.Page) *gorm.DB {
	return db.Offset(page.Offset()).Limit(page.Limit())
}

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		logger.Error(err)
	}
}
