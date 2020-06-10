package cmd

import (
	"fmt"
	"os"

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
	rootCmd.AddCommand(typedCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
