package main

import (
	l4g "code.google.com/p/log4go"
	"github.com/codegangsta/cli"
	"github.com/rocketjourney/swdaemon/model"
	"os"
	"time"
)

func init() {
	l4g.LoadConfiguration("config/log_level.xml")
}

func main() {

	messages := make(chan string)
	app := cli.NewApp()
	app.Name = "swdaemon"
	app.Usage = "Checkin Daemon"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "club",
			Value: "1",
			Usage: "club Identifier where daemon its running",
		},
	}
	app.Action = func(c *cli.Context) {
		println("Starting Daemon for club id: ", c.String("club"))
	}

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
