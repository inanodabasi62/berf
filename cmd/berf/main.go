package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/inanodabasi62/projects/berf/app"
)

func main() {
	rootCmd, _ := app.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", ""); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
