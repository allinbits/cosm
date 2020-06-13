package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches an application server.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		appName, _ := getAppAndModule()
		cmdInit := exec.Command("/bin/sh", "-c", "go mod tidy && make && sh init.sh")
		errInit := cmdInit.Run()
		if errInit != nil {
			fmt.Println(errInit.Error())
			return
		}
		cmdServer := exec.Command("/bin/sh", "-c", fmt.Sprintf("%[1]vd start", appName))
		errServer := cmdServer.Start()
		if errServer != nil {
			fmt.Println(errServer.Error())
			return
		}
		cmdREST := exec.Command(fmt.Sprintf("%[1]vcli", appName), "rest-server")
		errREST := cmdREST.Run()
		if errREST != nil {
			fmt.Println(errREST.Error())
			return
		}
	},
}
