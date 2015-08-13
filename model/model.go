package model

import (
	l4g "code.google.com/p/log4go"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jinzhu/gorm"
)

func InitDB() {
	db, err := gorm.Open("mssql", "server=swsql.ct2hxfttjshi.us-east-1.rds.amazonaws.com;user id=misael;password=00000000;database=master;encrypt=false")
	db.LogMode(true)

	if err != nil {
		l4g.Info(err)
	}
}
