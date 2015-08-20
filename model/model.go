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
	Net           network.Network
	Delay         int
	Query         string
	TimeFormat    string
}

func (m *Model) SetupModel() error {
	s := m.readSettings()
	l4g.Info(s)
	db, err := gorm.Open("mysql", s.User+":"+s.Password+"@tcp("+s.Server+":"+s.Port+")/"+s.DB_name+"?charset=utf8&parseTime=True&loc=Local")
	db.LogMode(true)
	m.DB = db
	m.Net = network.Network{}
	m.Net.Server = s.Path
	m.Delay = s.Delay
	m.Query = s.Query
	m.TimeFormat = s.Timeformat
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

	access := []Register{}
	const shortForm = "2006-01-02"
	searchDate := m.DateOfLastGet.Format(shortForm)
	searchHour := m.DateOfLastGet.Format(m.TimeFormat)

	m.DB.Where(m.Query, searchDate, searchHour).Find(&access)

	for _, r := range access {
		l4g.Info("%+v", r)
		m.Net.SendCheck(r.WayId, r.ClubId, r.UserId)
	}
	m.DateOfLastGet = time.Now()
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
