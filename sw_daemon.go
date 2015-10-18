package main

import (
	l4g "code.google.com/p/log4go"
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/inconshreveable/go-update"
	"github.com/rocketjourney/swdaemon/model"
	"github.com/rocketjourney/swdaemon/network"
	"os"
	"runtime"
	"time"
)

type JSONConfig struct {
	Changelog string                       `json:"changelog"`
	Version   string                       `json:"version"`
	Versions  map[string]map[string]string `json:"versions"`
	Pinned    map[string]string            `json:"pinned"`
}

func init() {
	l4g.LoadConfiguration("config/log_level.xml")
}

func main() {

	messages := make(chan string)
	app := cli.NewApp()
	app.Name = "swdaemon"
	app.Usage = "Checking Daemon"
	app.Version = model.VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "club",
			Value: "1",
			Usage: "club Identifier where daemon its running",
		},
	}
	app.Action = func(c *cli.Context) {
		l4g.Info("Running swdaemon version: %+v", model.VERSION)
		l4g.Info("Swdaemon proccess id: %+v", os.Getpid())
		l4g.Info("Starting Daemon for club id: ", c.String("club"))
	}

	l4g.Info("Checking for updates ...")
	checkForUpdate()

	model := model.Model{}
	err := model.SetupModel()
	if err == nil {
		go startSearch(&model)
	} else {
		l4g.Info(err)
	}

	app.Run(os.Args)
	msg := <-messages
	println("Stopping Daemon")
	println(msg)
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

func checkForUpdate() {
	l4g.Trace("Starting Updating process")
	net := network.Network{}
	file_data, err := net.GetUpdateFile(model.SERVER + model.UPDATE_PATH)

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
		err, _ := update.New().FromUrl(version[runtime.GOOS])
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

func isNewVersionIsAvailable(conf_data JSONConfig) bool {

	if model.VERSION != conf_data.Version {
		return true
	}
	return false
}
