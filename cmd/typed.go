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
	Use:   "type [typeName] [field1] [field2] ...",
	Short: "Generates CRUD actions for type",
	Args:  cobra.MinimumNArgs(1),
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
		var fields []typed.Field
		for _, f := range args[1:] {
			fs := strings.Split(f, ":")
			name := fs[0]
			var datatype string
			acceptedTypes := map[string]bool{
				"string": true,
				"bool":   true,
				"int":    true,
				"float":  true,
			}
			if len(fs) == 2 && acceptedTypes[fs[1]] {
				datatype = fs[1]
			} else {
				datatype = "string"
			}
			field := typed.Field{Name: name, Datatype: datatype}
			fields = append(fields, field)
		}
		g, _ := typed.New(&typed.Options{
			ModulePath: modulePath,
			AppName:    appName,
			TypeName:   args[0],
			Fields:     fields,
		})
		run := genny.WetRunner(context.Background())
		run.With(g)
		run.Run()
	},
}
