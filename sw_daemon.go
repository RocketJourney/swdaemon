package main

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/codegangsta/cli"
	"github.com/rocketjourney/swdaemon/model"
	"github.com/rocketjourney/swdaemon/routines"
	"os"
)

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
	model := model.Model{}
	dbConn := routines.DBConnection{}
	dbConn.Connect(&model)

	go func() {
		err := model.SetupModel()
		if err == nil {
			dbConn.Connected <- true
		} else {
			l4g.Info(err)
			dbConn.Connected <- false
		}
	}()

	app.Run(os.Args)
	msg := <-messages
	println("Stopping Daemon")
	println(msg)
}
