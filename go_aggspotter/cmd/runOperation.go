/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package spotter

import (
	"fmt"

	// log "github.com/sirupsen/logrus"
	"github.com/gagliardetto/solana-go"
	"github.com/spf13/cobra"
)

// runOperationCmd represents the runOperation command
var runOperationCmd = &cobra.Command{
	Use:   "runOperation",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("runOperation called")

		keeperAccsPaths, err := cmd.Flags().GetStringSlice("keeper")
		if err != nil {
			panic(err)
		}

		operation, err := cmd.Flags().GetString("operation")
		if err != nil {
			panic(err)
		}

		// Parse operation key
		operationPk := solana.MustPublicKeyFromBase58(operation)

		// Read keeper private keys
		keeperAccs := make([]solana.PrivateKey, len(keeperAccsPaths))
		for i, path := range keeperAccsPaths {
			keeperAccs[i], err = solana.PrivateKeyFromSolanaKeygenFile(path)
			if err != nil {
				panic(err)
			}
		}

		solanaConfigPath, err := cmd.Flags().GetString("keypair");
		if err != nil {
			panic(err)
		}

		client := NewSpotter(solanaConfigPath)
		client.ExecuteScript(keeperAccs, operationPk)
	},
}

func init() {
	rootCmd.AddCommand(runOperationCmd)

	runOperationCmd.Flags().StringSliceP("keeper", "k", []string{}, "List of keepers")
	runOperationCmd.Flags().SetAnnotation("keeper", cobra.BashCompOneRequiredFlag, []string{"true"})
	runOperationCmd.MarkFlagRequired("keeper")
	runOperationCmd.MarkFlagFilename("keeper", "json")

	runOperationCmd.Flags().StringP("operation", "o", "", "Operation public key")
	// runOperationCmd.Flags().String("keypair", os.Getenv("HOME") + "/.config/solana/id.json", "Operation public key")
	// runOperationCmd.MarkFlagFilename("keypair", "json")
}
