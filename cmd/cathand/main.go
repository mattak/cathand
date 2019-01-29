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
	recordingFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "size",
			Value: "720x1280",
			Usage: "screen size of recording video",
		},
		cli.IntFlag{
			Name:  "bit-rate",
			Value: 4000000,
			Usage: "bit rate of recording video",
		},
		cli.IntFlag{
			Name:  "time-limit",
			Value: 10,
			Usage: "time limit of each recording video. limit max is 180 sec",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "record",
			Aliases:   []string{"r"},
			Usage:     "record touch events",
			ArgsUsage: "[project_name]",
			Flags:     recordingFlags,
			Action: func(c *cli.Context) error {
				if len(c.Args()) <= 0 {
					return errors.New("ERROR: missing project name")
				}
				project := cathand.NewProject(c.Args().First(), "")
				option := cathand.NewRecordOption()
				option.Size = c.String("size")
				option.BitRate = c.Uint("bit-rate")
				option.TimeLimit = c.Uint("time-limit")

				cathand.CommandRecord(project, option)
				return nil
			},
		},
		{
			Name:      "compose",
			Aliases:   []string{"c"},
			Usage:     "compose playable touch events from recorded data",
			ArgsUsage: "[record_project_name] [play_project_name]",
			Action: func(c *cli.Context) error {
				if len(c.Args()) <= 1 {
					return errors.New("ERROR: requires 2 project name: [recorded_project_name] [play_project_name]")
				}
				recordProject := cathand.NewProject(c.Args().Get(0), "")
				playProject := cathand.NewProject(c.Args().Get(1), "")
				cathand.CommandCompose(recordProject, playProject)
				return nil
			},
		},
		{
			Name:      "play",
			Aliases:   []string{"p"},
			Usage:     "play playable touch events",
			ArgsUsage: "[play_project_name] [result_project_name]",
			Flags:     recordingFlags,
			Action: func(c *cli.Context) error {
				if len(c.Args()) <= 1 {
					return errors.New("ERROR: requires 2 project name: [recorded_project_name] [play_project_name]")
				}
				playProject := cathand.NewProject(c.Args().Get(0), "")
				resultProject := cathand.NewProject(c.Args().Get(1), "")

				option := cathand.NewRecordOption()
				option.Size = c.String("size")
				option.BitRate = c.Uint("bit-rate")
				option.TimeLimit = c.Uint("time-limit")

				cathand.CommandPlay(playProject, resultProject, option)
				return nil
			},
		},
		{
			Name:      "split",
			Aliases:   []string{"s"},
			Usage:     "split video into image segments",
			ArgsUsage: "[project_name]+",
			Action: func(c *cli.Context) error {
				if len(c.Args()) <= 1 {
					return errors.New("ERROR: missing project name")
				}

				projects := make([]cathand.Project, len(c.Args()))
				for i := 0; i < len(c.Args()); i++ {
					projects[i] = cathand.NewProject(c.Args().Get(i), "")
				}
				cathand.CommandSplit(projects...)
				return nil
			},
		},
		{
			Name:      "verify",
			Aliases:   []string{"v"},
			Usage:     "verify recorded project and result project",
			ArgsUsage: "[result_project_name] [record_project_name]",
			Action: func(c *cli.Context) error {
				if len(c.Args()) <= 1 {
					return errors.New("ERROR: requires 2 project name: [recorded_project_name] [result_project_name]")
				}
				project1 := cathand.NewProject(c.Args().Get(0), "")
				project2 := cathand.NewProject(c.Args().Get(1), "")
				cathand.CommandVerify(project1, project2)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
