package main

import (
	"github.com/ngaut/log"
	"github.com/spf13/cobra"
)

func main() {
	log.Info("Run jepsenFuzz")
	var rootCmd = &cobra.Command{
		Use:   "jepsenFuzz",
		Short: "jepsenFuzz toolset",
	}
	rootCmd.AddCommand(newInitCmd())
	rootCmd.Execute()

}
