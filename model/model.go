package model

import (
	l4g "code.google.com/p/log4go"
	"encoding/json"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/rocketjourney/swdaemon/network"
	"io/ioutil"
	"time"
)

type Model struct {
	DB            gorm.DB
	DateOfLastGet time.Time
}

func (m *Model) SetupModel() error {
	s := m.readSettings()
	db, err := gorm.Open("mysql", s.User+":"+s.Password+"@tcp("+s.Server+":"+s.Port+")/"+s.DB_name+"?charset=utf8&parseTime=True&loc=Local")
	db.LogMode(true)
	m.DB = db
	m.DateOfLastGet = time.Now()

	if err != nil {
		l4g.Info(err)
		return err
	} else {
		ping_err := db.DB().Ping()
		if ping_err != nil {
			l4g.Info(ping_err)
			return ping_err
		}
	}
	return nil
}

func (m *Model) SearchAccess() {
	acceso := Register{}
	m.DB.First(&acceso)
	l4g.Info("%+v", acceso)
	println("check-in")
	network.SendCheck(true)
	m.DateOfLastGet = time.Now()
	l4g.Info("%+v", m.DateOfLastGet)
}

func (m *Model) readSettings() *Settings {
	dat, _ := ioutil.ReadFile("config.json")
	settings := Settings{}
	err := json.Unmarshal(dat, &settings)
	if err != nil {
		l4g.Info("error:", err)
	}
	l4g.Info("%+v", settings)
	return &settings
}
