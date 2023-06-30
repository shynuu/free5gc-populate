package main

import (
	"log"
	"os"

	"github.com/shynuu/free5gc-populate/runtime"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetPrefix("[free5gc-populate]")
	var config string = ""

	app := &cli.App{
		Name:                 "free5gc-populate",
		Usage:                "Populate free5gc mondo database from a config file",
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			{Name: "Youssouf Drif"},
		},
		Copyright: "Copyright (c) 2021 Youssouf Drif",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Load configuration from `FILE`",
				Destination: &config,
				Required:    true,
				DefaultText: "not set",
			},
		},
		Action: func(c *cli.Context) error {
			err := runtime.Run(config)
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
