package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"strings"

	"github.com/allinbits/cosmos-cli/templates/typed"
	"github.com/gobuffalo/genny"
	"github.com/spf13/cobra"
)

var typedCmd = &cobra.Command{
	Use:   "type [typeName]",
	Short: "Generates CRUD actions for type",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		goModFile, err := ioutil.ReadFile("go.mod")
		if err != nil {
			log.Fatal(err)
		}
		moduleString := strings.Split(string(goModFile), "\n")[0]
		modulePath := strings.ReplaceAll(moduleString, "module ", "")
		var appName string
		if t := strings.Split(modulePath, "/"); len(t) > 0 {
			appName = t[len(t)-1]
		}
		g, _ := typed.New(&typed.Options{
			ModulePath: modulePath,
			AppName:    appName,
			TypeName:   args[0],
		})
		run := genny.WetRunner(context.Background())
		run.With(g)
		run.Run()
	},
}
