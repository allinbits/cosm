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

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Launches an application server.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		appName, _ := getAppAndModule()
		fmt.Printf("\n📦 Installing dependencies...\n")
		cmdMod := exec.Command("/bin/sh", "-c", "go mod tidy")
		cmdMod.Stdout = os.Stdout
		if err := cmdMod.Run(); err != nil {
			log.Fatal("Error running go mod tidy. Please, check ./go.mod")
		}
		fmt.Printf("🚧 Building the application...\n")
		cmdMake := exec.Command("/bin/sh", "-c", "make")
		cmdMake.Stdout = os.Stdout
		if err := cmdMake.Run(); err != nil {
			log.Fatal("Error in building the application. Please, check ./Makefile")
		}
		fmt.Printf("💫 Initializing the chain...\n")
		cmdInit := exec.Command("/bin/sh", "-c", "sh init.sh")
		cmdInit.Stdout = os.Stdout
		if err := cmdInit.Run(); err != nil {
			log.Fatal("Error in initializing the chain. Please, check ./init.sh")
		}
		fmt.Printf("🌍 Running a server at http://localhost:26657 (Tendermint)\n")
		cmdTendermint := exec.Command("/bin/sh", "-c", fmt.Sprintf("%[1]vd start", appName))
		cmdTendermint.Stdout = os.Stdout
		if err := cmdTendermint.Start(); err != nil {
			log.Fatal(fmt.Sprintf("Error in running %[1]vd start", appName))
		}
		fmt.Printf("🌍 Running a server at http://localhost:1317 (LCD)\n\n")
		cmdREST := exec.Command(fmt.Sprintf("%[1]vcli", appName), "rest-server")
		cmdREST.Stdout = os.Stdout
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

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches a hot-reloading application server.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("📦 Making sure hot reloading is enabled...\n")
		cmdAirGet := exec.Command("/bin/sh", "-c", "GO111MODULE=off go get github.com/cosmtrek/air")
		if err := cmdAirGet.Run(); err != nil {
			log.Fatal("Error in enabling hot reload with air.")
		}
		cmdAir := exec.Command("/bin/sh", "-c", "air -d")
		cmdAir.Stdout = os.Stdout
		if err := cmdAir.Run(); err != nil {
			log.Fatal("Error in running `cosm serve` with `air`.")
		}
	},
}
