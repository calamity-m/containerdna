package cmd

import (
	"fmt"
	"github.com/calamity-m/containerdna/pkg/heritage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	Child string

	Parents []string

	strictCmd = &cobra.Command{
		Use:     "strict",
		Short:   "Do a strict heritage check",
		GroupID: heritageGroup.ID,
		Long: `Complete a strict heritage check, which verifies that for every parent provided the child must originate
from every single one.

This is done on a layer comparison basis. Given an example parent1 and parent2, the child must contain all layers of a
specified parent from its initial layer

	parent1 	-> layer0: A
	
	parent2 	-> layer0: A
		   		-> layer1: AA

	child1 		-> layer0: A
				-> layer1: AA
				-> layer2: AAA

	child2      -> layer0: A

	Child 1 is built from parent1 and parent2
	Child 2 is not built from parent1 and parent2, as it is lacking parent2's second layer.

Usage:
	containerdna --child docker://nginx --parent docker://nginx --parent docker://nginx --parent docker://nginx
`,
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Debug("Entering strictCmd")
			if len(Parents) == 0 || Child == "" {
				fmt.Printf("Please provide child and parent\\s\n")
			}

			childRef, err := heritage.GetImageReference(Child)
			if err != nil {
				os.Exit(1)
			}
			// This should be in a goroutine
			childLayers, err := heritage.GetImageLayers(childRef)
			if err != nil {
				os.Exit(1)
			}

			parentRef, err := heritage.GetImageReference(Parents[0])
			if err != nil {
				os.Exit(1)
			}
			// This should be in a goroutine
			parentLayers, err := heritage.GetImageLayers(parentRef)
			if err != nil {
				os.Exit(1)
			}

			v := heritage.ValidateChildParents(childLayers, parentLayers)

			fmt.Printf("Valid Child: %v\n", v)
			if !v {
				os.Exit(1)
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(strictCmd)

	strictCmd.Flags().StringSliceVarP(&Parents, "parent", "p", []string{}, "Supposed Parent Image")
	strictCmd.Flags().StringVarP(&Child, "child", "c", "", "Supposed Child Image")
	strictCmd.MarkFlagRequired("parent")
	strictCmd.MarkFlagRequired("child")
}

func Strict(child string, parents ...string) (bool, error) {
	// Need to implement

	return false, nil
}
