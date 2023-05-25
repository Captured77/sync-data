package main

import (
	"fmt"
	"log"
	"os"
	"github.com/urfave/cli"
)

const (
	downloadUrl = "http://10.113.3.1/swapp/release/deeplearningsamples/"
	fileName    = "download_deps.sh"
	defaultPath = "/home/huoshangbing/workspace"
)

func main() {
	var url string
	var file string
	var duration int
	var path string

	app := cli.NewApp()
	app.Name = "sync data"
	app.Usage = "sync data from url"
	app.Version = "v0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "url, u",
			Usage:       "download shell script url",
			Value:       downloadUrl,
			Destination: &url,
		},
		cli.StringFlag{
			Name:        "file, f",
			Usage:       "shell script",
			Value:       fileName,
			Destination: &file,
		},
		cli.IntFlag{
			Name:        "duration, d",
			Usage:       "sync duration (day)",
			Value:       7,
			Destination: &duration,
		},
		cli.StringFlag{
			Name:        "path, p",
			Usage:       "download data path, if not exits, program will create!",
			Value:       defaultPath,
			Destination: &path,
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println(url)
		fmt.Println(file)
		fmt.Println(duration)
		fmt.Println(path)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal()
	}
	
	
}
