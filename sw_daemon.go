package main

import (
	"github.com/codegangsta/cli"
	"github.com/rocketjourney/swdaemon/model"
	"os"
	"time"
)

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

	go checkIn()

	app.Run(os.Args)
	msg := <-messages
	println("Stopping Daemon")
	println(msg)
}

func checkIn() {
	model.InitDB()
	for {
		println("check-in")
		delay := (time.Second * time.Duration(1))
		time.Sleep(delay)
	}
}
