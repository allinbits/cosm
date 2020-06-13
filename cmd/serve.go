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
		fmt.Printf("\n📦 Installing dependencies...\n")
		cmdMod := exec.Command("/bin/sh", "-c", "go mod tidy")
		errMod := cmdMod.Run()
		if errMod != nil {
			fmt.Println(errMod.Error())
			return
		}
		fmt.Printf("🚧 Building the application...\n")
		cmdMake := exec.Command("/bin/sh", "-c", "make")
		errMake := cmdMake.Run()
		if errMake != nil {
			fmt.Println(errMake.Error())
			return
		}
		fmt.Printf("💫 Initializing the chain...\n")
		cmdInit := exec.Command("/bin/sh", "-c", "sh init.sh")
		errInit := cmdInit.Run()
		if errInit != nil {
			fmt.Println(errInit.Error())
			return
		}
		fmt.Printf("🌍 Running a server at http://localhost:26657 (Tendermint)\n")
		cmdServer := exec.Command("/bin/sh", "-c", fmt.Sprintf("%[1]vd start", appName))
		errServer := cmdServer.Start()
		if errServer != nil {
			fmt.Println(errServer.Error())
			return
		}
		fmt.Printf("🌍 Running a server at http://localhost:1317 (LCD)\n\n")
		cmdREST := exec.Command(fmt.Sprintf("%[1]vcli", appName), "rest-server")
		errREST := cmdREST.Run()
		if errREST != nil {
			fmt.Println(errREST.Error())
			return
		}
	},
}
