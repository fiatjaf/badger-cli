package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version: "1.0.0",
	Use:     "badger-cli",
	Short:   "Command line client for managing a badger database",
}

func main() {
	rootCmd.PersistentFlags().StringP("dir", "d", "", "Path to the badger database direcotry")
	rootCmd.MarkFlagRequired("dir")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
