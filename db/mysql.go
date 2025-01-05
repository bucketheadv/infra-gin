package db

import (
	"database/sql"
	"github.com/bucketheadv/infragin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySqlConf struct {
	Url string
}

func NewMySQL(config MySqlConf, gormConfig *gorm.Config) *gorm.DB {
	var err error
	DB, err := gorm.Open(mysql.Open(config.Url), gormConfig)
	if err != nil {
		logrus.Fatalln(err)
	}
	return DB
}

func Page(db *gorm.DB, page infragin.Page) *gorm.DB {
	return db.Offset(page.Offset()).Limit(page.Limit())
}

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		logrus.Error(err)
	}
}
