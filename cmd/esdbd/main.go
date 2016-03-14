package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/cli"

	"github.com/eventsourcedb/eventsourcedb"
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eventsourcedb"
	app.Usage = "inflate balloon"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "cert, c",
			Value:  "./cert.pem",
			Usage:  "path to TLS certificate file",
			EnvVar: "BALLOON_CERT",
		},
		cli.StringFlag{
			Name:   "key, k",
			Value:  "./key.pem",
			Usage:  "path to TLS key file",
			EnvVar: "BALLOON_KEY",
		},
		cli.StringFlag{
			Name:   "data, d",
			Usage:  "Path to datadir",
			EnvVar: "BALLOON_DATA",
		},
		cli.StringFlag{
			Name:   "addr, i",
			Value:  ":8080",
			Usage:  "Address to bind to",
			EnvVar: "BALLOON_ADDR",
		},
	}
	app.Action = MainCli
	return app
}

func MainCli(c *cli.Context) {
	mux := eventsourcedb.Routes()
	addr := c.String("addr")
	err := http.ListenAndServeTLS(addr, c.String("cert"), c.String("key"), mux)
	log.Fatal(err)
}

func main() {
	app := NewApp()
	app.RunAndExitOnError()
}
