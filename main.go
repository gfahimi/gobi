package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anuvu/gobi/builder"
	"github.com/urfave/cli"
)

var (
	version    = ""
	buildDir   = "build"
	configFile = "build.yaml"
)

func main() {
	app := cli.NewApp()
	app.Name = "arnold"
	app.Usage = "Arnold is the builder for Go"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config-file, f",
			Usage: "load configuration from `FILE`",
			Value: "build.yaml",
		},
		cli.StringFlag{
			Name:  "output-dir, o",
			Usage: "set the directory for build output relative to project root",
			Value: "./build",
		},
		cli.StringFlag{
			Name:  "project-dir, p",
			Usage: "set the directory for go project root",
			Value: ".",
		},
	}

	// Default action is to build all
	app.Action = builder.Build

	app.Before = func(ctx *cli.Context) error {
		var err error
		if buildDir, err = filepath.Abs(ctx.String("output-dir")); err != nil {
			return err
		}
		if configFile, err = filepath.Abs(ctx.String("config-file")); err != nil {
			return err
		}

		if e := os.RemoveAll(buildDir); e != nil {
			return e
		}
		return os.Mkdir(buildDir, os.ModeDir|os.ModePerm)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Println("=== Build Failed ===")
		os.Exit(1)
	}

	fmt.Println("=== Build succeeded ===")
}
