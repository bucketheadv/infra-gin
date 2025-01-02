package db

import (
	"database/sql"
	"github.com/bucketheadv/infragin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Page(db *gorm.DB, page infragin.Page) *gorm.DB {
	return db.Offset(page.Offset()).Limit(page.Limit())
}

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		logrus.Error(err)
	}
}
