package main

import (
	"os"

	"github.com/cvzamannow/E-Learning-API/cmd"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warnln("[main] Error loading .env file")
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()

	app.Name = "Simaku E-Learning-API"

	app.Commands = []cli.Command{
		cmd.HTTPGatewayServerCMD(),
		cmd.DoMigrateUpCMD(),
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("[main] Failed running command becuase: %s", err.Error())
		os.Exit(1)
	}
}
