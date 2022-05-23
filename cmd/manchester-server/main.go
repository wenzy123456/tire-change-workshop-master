package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/surmus/tire-change-workshop/internal/manchester"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"time"
)

const (
	version        = "v2.0.0"
	listenPortFlag = "port"
	verboseFlag    = "verbose"
	defaultPort    = 9004
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    listenPortFlag,
		Aliases: []string{"p"},
		Value:   fmt.Sprintf("%d", defaultPort),
		Usage:   "Port for server to listen incoming connections",
	},
	&cli.BoolFlag{
		Name:  verboseFlag,
		Usage: "Enables debug messages print with SQL logging",
	},
}

// @title Manchester tire workshop API
// @version 2.0
// @description Tire workshop service IOT integration.
// @BasePath /api/v2
// @license.name MIT
func main() {
	app := cli.NewApp()
	app.Version = version
	app.Usage = "Manchester tire workshop API server"
	app.Flags = flags
	app.Action = initServer

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func initServer(c *cli.Context) error {
	listenToPort := c.Uint(listenPortFlag)

	if listenToPort == 0 {
		return fmt.Errorf("invalid server listen port supplied: %s", c.String(listenPortFlag))
	}

	if c.Bool(verboseFlag) {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	return setupServer(listenToPort, c.Bool(verboseFlag))
}

func setupServer(port uint, debugMode bool) error {
	apiRouter := manchester.Init(debugMode)
	// The url pointing to API definition
	swaggerURL := ginSwagger.URL("swagger/doc.json")
	apiRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	workshopServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      apiRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Infof("application initialized, listening to port %d", port)
	return workshopServer.ListenAndServe()
}
