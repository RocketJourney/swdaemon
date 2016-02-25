package routines

import (
	l4g "github.com/alecthomas/log4go"
	"encoding/json"
	"fmt"
	"github.com/inconshreveable/go-update"
	"github.com/rocketjourney/swdaemon/model"
	"github.com/rocketjourney/swdaemon/network"
	"net/http"
	"os"
	"runtime"
	"time"
)

type DBConnection struct {
	Connected chan bool
	model     *model.Model
}

type JSONConfig struct {
	Changelog string                       `json:"changelog"`
	Version   string                       `json:"version"`
	Versions  map[string]map[string]string `json:"versions"`
	Pinned    map[string]string            `json:"pinned"`
}

func (dbConn *DBConnection) Connect(m *model.Model) {
	dbConn.Connected = make(chan bool)
	dbConn.model = m
	go dbConn.performConection()
}

func (dbConn *DBConnection) performConection() {
	for {
		select {
		case conn := <-dbConn.Connected:
			l4g.Info(conn)
			if conn {
				l4g.Info("Connected")
				l4g.Info("Checking for updates ...")
				checkForUpdate()
				go startSearch(dbConn.model)
				go reportAlive(dbConn.model)
				go scheduledUpdate()
			} else {
				l4g.Info("Retry Connection")
				go dbConn.retryConnection()
			}
		}
	}
}

func (dbConn *DBConnection) retryConnection() {
	delay := (time.Second * 5)
	time.Sleep(delay)
	err := dbConn.model.SetupModel()
	if err == nil {
		dbConn.Connected <- true
	} else {
		l4g.Info(err)
		dbConn.Connected <- false
	}
}

func startSearch(m *model.Model) {
	for {
		date := time.Now()
		if date.Hour() >= m.StandByStartHour && date.Hour() < m.StandByEndHour {

		} else {
			m.SearchAccess()
			delay := (time.Second * time.Duration(m.Delay))
			time.Sleep(delay)
		}
	}
}

func scheduledUpdate() {
	for {
		delay := (time.Minute * 2)
		time.Sleep(delay)
		checkForUpdate()
	}
}

func reportAlive(m *model.Model) {
	net := network.Network{}
	s := m.ReadSettings()
	l4g.Info(s)
	net.Server = model.SERVER
	net.AccessToken = s.Access_token
	pid := fmt.Sprintf("%d", os.Getpid())
	club := fmt.Sprintf("%d", s.Spot_id)
	for {
		net.ReportAlive(pid, club)
		delay := (time.Minute * 10)
		time.Sleep(delay)
	}
}

func checkForUpdate() {
	l4g.Trace("Starting Updating process")
	net := network.Network{}
	file_data, err := net.GetUpdateFile(model.UPDATE_SERVER + model.UPDATE_PATH)

	if err != nil {
		l4g.Trace("Error getting update file")
		return
	}

	conf_data := JSONConfig{}
	if err := json.Unmarshal(*file_data, &conf_data); err != nil {
		l4g.Error("Error parsing config file", err)
	}

	if isNewVersionIsAvailable(conf_data) {
		l4g.Info("New SW version available. Downloading: %s", conf_data.Version)
		version := conf_data.Versions[conf_data.Version]
		err := doUpdate(version[runtime.GOOS])
		if err != nil {
			l4g.Error("Update failed: %v", err)
			return
		} else {
			l4g.Info("New version ready. Please restart swdaemon to using it.")
			l4g.Info("SW Daemon version %s changelog:", conf_data.Version)
			l4g.Info("%s", conf_data.Changelog)
		}
	} else {
		l4g.Info("No new updates founded")
	}
}

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// error handling
	}
	return err
}

func isNewVersionIsAvailable(conf_data JSONConfig) bool {

	if model.VERSION != conf_data.Version {
		return true
	}
	return false
}
