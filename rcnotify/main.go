package main

import (
	"log"
	"os"
	"time"

	"github.com/bradgignac/cloud-notifications/config"
	"github.com/bradgignac/cloud-notifications/ingestor"
	"github.com/bradgignac/cloud-notifications/notifier"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rcnotify"
	app.Usage = "Push notifications for Rackspace Cloud"
	app.Version = "0.0.0"
	app.Author = "Brad Gignac"
	app.Email = "bgignac@bradgignac.com"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "rackspace-user",
			Usage: "Rackspace Cloud username",
		},
		cli.StringFlag{
			Name:  "rackspace-key",
			Usage: "Rackspace Cloud API key",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Cloud Notifications config file",
		},
	}

	app.Action = poll

	app.Run(os.Args)
}

func poll(c *cli.Context) {
	config, err := config.LoadYAML(c.String("config"))
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
		os.Exit(1)
	}

	notifier, err := notifier.New(config.Notifier.Type, config.Notifier.Options)
	if err != nil {
		log.Fatalf("Failed to instantiate notifier: %v", err)
		os.Exit(1)
	}

	rsUser := arg(c, "rackspace-user")
	rsKey := arg(c, "rackspace-key")

	ingestor := &ingestor.CloudFeeds{
		Notifier: notifier,
		Interval: 10 * time.Second,
		User:     rsUser,
		Key:      rsKey,
	}

	err = ingestor.Start()
	if err != nil {
		log.Fatalf("Failed to start ingestor: %v", err)
		os.Exit(1)
	}
}

func arg(c *cli.Context, name string) string {
	val := c.String(name)
	if val == "" {
		log.Fatalf("Parameter \"%s\" was not provided\n", name)
	}

	return val
}
