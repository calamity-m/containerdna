package cmd

import (
	"fmt"
	"github.com/calamity-m/containerdna/pkg/heritage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	Parent string

	Child string

	parents []string

	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Debug("Entering checkCmd")
			if Parent == "" || Child == "" {
				return
			}
			parentRef, err := heritage.GetImageReference(Parent)
			if err != nil {
				os.Exit(1)
			}
			childRef, err := heritage.GetImageReference(Child)
			if err != nil {
				os.Exit(1)
			}

			// This should be in a goroutine
			parentLayers, err := heritage.GetImageLayers(parentRef)
			if err != nil {
				os.Exit(1)
			}

			// This should be in a goroutine
			childLayers, err := heritage.GetImageLayers(childRef)
			if err != nil {
				os.Exit(1)
			}

			v := heritage.ValidateChildParents(childLayers, parentLayers)

			fmt.Printf("Valid Child: %v\n", v)
		},
	}
)

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringSliceVar(&parents, "parents", []string{}, "Parent")
	checkCmd.Flags().StringVarP(&Parent, "parent", "p", "", "Supposed Parent Image")
	checkCmd.Flags().StringVarP(&Child, "child", "c", "", "Supposed Child Image")
	checkCmd.MarkFlagRequired("parent")
	checkCmd.MarkFlagRequired("child")
}
