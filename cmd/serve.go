package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func setupCloseHandler(tendermint *exec.Cmd, rest *exec.Cmd) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := tendermint.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
		if err := rest.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
		os.Exit(0)
	}()
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches an application server.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		appName, _ := getAppAndModule()
		fmt.Printf("\n📦 Installing dependencies...\n")
		cmdMod := exec.Command("/bin/sh", "-c", "go mod tidy")
		if err := cmdMod.Run(); err != nil {
			log.Fatal("Error running go mod tidy. Please, check ./go.mod")
		}
		fmt.Printf("🚧 Building the application...\n")
		cmdMake := exec.Command("/bin/sh", "-c", "make")
		if err := cmdMake.Run(); err != nil {
			log.Fatal("Error in building the application. Please, check ./Makefile")
		}
		fmt.Printf("💫 Initializing the chain...\n")
		cmdInit := exec.Command("/bin/sh", "-c", "sh init.sh")
		if err := cmdInit.Run(); err != nil {
			log.Fatal("Error in initializing the chain. Please, check ./init.sh")
		}
		fmt.Printf("🌍 Running a server at http://localhost:26657 (Tendermint)\n")
		cmdTendermint := exec.Command("/bin/sh", "-c", fmt.Sprintf("%[1]vd start", appName))
		if err := cmdTendermint.Start(); err != nil {
			log.Fatal(fmt.Sprintf("Error in running %[1]vd start", appName))
		}
		fmt.Printf("🌍 Running a server at http://localhost:1317 (LCD)\n\n")
		cmdREST := exec.Command(fmt.Sprintf("%[1]vcli", appName), "rest-server")
		if err := cmdREST.Start(); err != nil {
			log.Fatal(fmt.Sprintf("Error in running %[1]vcli rest-server", appName))
		}
		fmt.Printf("🔧 Running dev interface at http://localhost:12345\n\n")
		setupCloseHandler(cmdTendermint, cmdREST)
		router := mux.NewRouter()
		cosmUI := packr.New("ui/dist", "../ui/dist")
		router.PathPrefix("/").Handler(http.FileServer(cosmUI))
		log.Fatal(http.ListenAndServe(":12345", router))
	},
}
