package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	dnaCmd = &cobra.Command{
		Use:     "dna",
		Short:   "Get history information from image",
		GroupID: containerGroup.ID,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dna")
		},
	}
)

func init() {
	rootCmd.AddCommand(dnaCmd)
}
