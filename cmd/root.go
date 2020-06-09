package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/allinbits/cosmos-cli/genny/app"
	"github.com/gobuffalo/genny"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cosmos",
	Short: "A tool for scaffolding out Cosmos applications",
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(appCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var appCmd = &cobra.Command{
	Use:   "app [github.com/org/repo]",
	Short: "Generates an empty application boilerplate",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var appName string
		if t := strings.Split(args[0], "/"); len(t) > 0 {
			appName = t[len(t)-1]
		}
		g, _ := app.New(&app.Options{
			ModulePath: args[0],
			AppName:    appName,
		})
		run := genny.WetRunner(context.Background())
		run.With(g)
		pwd, _ := os.Getwd()
		run.Root = pwd + "/" + appName
		run.Run()
	},
}
