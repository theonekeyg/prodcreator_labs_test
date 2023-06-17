/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package spotter

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go_aggspotter",
	Short: "CLI to interact with AggregationSpotter contract",
	Long: `CLI to interact with AggregationSpotter on-chain contract.
Currently supports only solana blockchain.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	runOperationCmd.Flags().StringP("url", "u", "", "RPC endpoint to solana node")
}
