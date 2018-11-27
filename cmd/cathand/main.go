package main

import (
	"errors"
	"github.com/mattak/cathand/pkg/cathand"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "cathand"
	app.Usage = "record then play just like you on android device"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:      "record",
			Aliases:   []string{"r"},
			Usage:     "record touch events",
			ArgsUsage: "[project_name]",
			Action: func(c *cli.Context) error {
				if len(c.Args()) <= 0 {
					return errors.New("ERROR: missing project name")
				}
				cathand.CommandRecord(c.Args().First())
				return nil
			},
		},
		{
			Name:      "compose",
			Aliases:   []string{"c"},
			Usage:     "compose playable touch events from recorded data",
			ArgsUsage: "[project_name]",
			Action: func(c *cli.Context) error {
				if len(c.Args()) <= 0 {
					return errors.New("ERROR: missing project name")
				}
				cathand.CommandCompose(c.Args().First())
				return nil
			},
		},
		{
			Name:      "play",
			Aliases:   []string{"p"},
			Usage:     "play playable touch events",
			ArgsUsage: "[project_name]",
			Action: func(c *cli.Context) error {
				if len(c.Args()) <= 0 {
					return errors.New("ERROR: missing project name")
				}
				cathand.CommandPlay(c.Args().First())
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
