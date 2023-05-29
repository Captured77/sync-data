package main

import (
	"strings"
	"bufio"
	"fmt"
	"log"
	"io"
	"net/http"
	neturl "net/url"
	"os"

	"github.com/urfave/cli"
	"os/exec"
	"github.com/jasonlvhit/gocron"
	"path/filepath"
)

const (
	downloadUrl = "http://10.113.3.1/swapp/release/deeplearningsamples/"
	fileName    = "download_deps.sh"
	defaultPath = "/home/ubuntu/huoshangbing/workspace"
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
			Usage:       "sync duration (day),default 7 days;Don't rather than 30 days",
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
		// url 校验
		_, err := neturl.ParseRequestURI(url)
		if err != nil {
			log.Fatal(err)
		}
		// duration
		if duration > 30 {
			log.Fatal("duration don't rather than 30")
		}
		// path check, if not exists, create 
		if _, err = os.Stat(path);  err != nil {
			err := os.MkdirAll(path, 0777); if err !=nil {
				log.Fatal(err)
			}
		}


		var shellCmd []string = rawDownloadScript(downloadUrl, file, path, "download.sh")
		var dataVersion = strings.Split(shellCmd[0], "\"")[1]
		var dataUrl = strings.Replace(shellCmd[1], "${version}", dataVersion, -1)
		var untarPackage = strings.Split(strings.Replace(shellCmd[3], "${version}", dataVersion, -1), " ")

		// download tar package
		if err = downloadDatasource(dataUrl, path); err != nil {
			log.Fatal(err)
		}

		//execute cmd, unzip .tar
		if err = untarDatasource(untarPackage, path); err != nil {
			log.Fatal(err)
		}

		s := gocron.NewScheduler()
		s.Every(uint64(duration)).Day().Do(func() {
			var shellCmdNew []string = rawDownloadScript(downloadUrl, file, path, "new-download.sh")
			var dataVersionNew = strings.Split(shellCmdNew[0], "\"")[1]
			if dataVersion != dataVersionNew {
				err := os.Remove(fmt.Sprintf("%s/%s", path, "download.h"))
				if err != nil {
					log.Fatal(err)
				}
				err = os.Rename(fmt.Sprintf("%s/%s", path, "new-download.sh"), fmt.Sprintf("%s/%s", path, "download.sh"))
				if err != nil {
					log.Fatal(err)
				}

				// delelte old core-*
				files, err := filepath.Glob("path/core-*")
				if err != nil {
					log.Fatal(err)
				}
				for _, f := range files {
					if err := os.RemoveAll(f); err != nil {
						log.Fatal(err)
					}
				}

				// delete data path
				if err := os.RemoveAll(fmt.Sprintf("%s/%s", path, "data")); err != nil {
					log.Fatal(err)
				}

				var dataVersion = strings.Split(shellCmdNew[0], "\"")[1]
				var dataUrl = strings.Replace(shellCmdNew[1], "${version}", dataVersion, -1)
				var untarPackage = strings.Split(strings.Replace(shellCmdNew[3], "${version}", dataVersion, -1), " ") 
				
				// download new data tar
				if err := downloadDatasource(dataUrl, path); err != nil {
					log.Fatal(err)
				}

				// untar package
				if err := untarDatasource(untarPackage, path); err != nil {
					log.Fatal(err)
				}

			}
			err := os.Remove(fmt.Sprintf("%s/%s", path, "new-download.sh"))
			if err != nil {
				log.Fatal(err)
			}
		})
		<-s.Start()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func rawDownloadScript(downloadurl, file, path, filename string) []string {
	// get download shell script 
	var shellCmd []string
	resp, err := http.Get(fmt.Sprintf("%s%s", downloadUrl, file))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	
	// create download shell
	out, err := os.Create(fmt.Sprintf("%s/%s", path, filename)); if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	f, err := os.Open(fmt.Sprintf("%s/%s", path, "download.sh"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		shellCmd = append(shellCmd, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return shellCmd
}

func downloadDatasource(dataurl, path string) error {
	//execute cmd, download datasource
	cmd := exec.Command(strings.Split(dataurl, " ")[0], strings.Split(dataurl, " ")[1], "-P", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func untarDatasource(untarparam []string, path string) error {
	cmd := exec.Command(untarparam[0], untarparam[1], fmt.Sprintf("%s/%s", path, untarparam[2]), "-C", path)
	cmd.Stdout =  os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}