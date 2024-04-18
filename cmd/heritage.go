package cmd

import (
	"fmt"
	"os"

	"github.com/calamity-m/containerdna/pkg/heritage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Relaxed bool

	Child string

	Parents []string

	heritageCmd = &cobra.Command{
		Use:     "heritage",
		Short:   "Do a heritage check",
		GroupID: containerGroup.ID,
		Long: `Complete a heritage check, which verifies that for every parent provided the child must originate
from every single one.

This is done on a layer comparison basis. Given an example parent1 and parent2, the child must contain all layers of a
specified parent from its initial layer

	parent1 - -> layer0: A
	
	parent2 - -> layer0: A
		  -> layer1: AA

	child1  - -> layer0: A
		  -> layer1: AA
		  -> layer2: AAA

	child2  - -> layer0: A

With the default strict check:

	Child 1 is built from parent1 and parent2
	Child 2 is not built from parent1 and parent2, as it is lacking parent2's second layer.

With --relaxed

	Child 1 and Child 2 are valid, as at least one parent is in their history

Usage:

	Strict and relaxed checks:

	containerdna --child docker://nginx --parent docker://nginx --parent docker://nginx --parent docker://nginx
	containerdna --relaxed --child docker://nginx --parent docker://alpine --parent docker://nginx

	Running against local images:

	containerdna --child docker-daemon:alpine:latest --parent docker-daemon:alpine:latest

`,
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Debug("Entering heritageCmd")
			if len(Parents) == 0 || Child == "" {
				fmt.Printf("Please provide child and parent\\s\n")
			}

			// heritage.Playground()
			//heritage.ValidateWithChannelsNoWg(false, "docker://alpine", "docker://nginx", "docker://ubuntu", "docker://alpine")
			valid, err := heritage.GetImageAsyncChannels(false, "docker://alpine", "docker://nginxasdfasdf", "adocker://ubuntu", "docker://alpine")
			if err != nil {
				fmt.Println("Encountered errors while attempting to evaluate heritage. Errors are:")
				fmt.Println(err)
				os.Exit(1)
			} else {
				fmt.Printf("Validity: %t", valid)
			}
			/*
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

				// v := heritage.ValidateChildParents(childLayers, parentLayers)

				fmt.Printf("Valid Child: %v\n", v)
				if !v {
					os.Exit(1)
				}
			*/

		},
	}
)

func init() {
	rootCmd.AddCommand(heritageCmd)

	heritageCmd.Flags().StringVarP(&Child, "child", "c", "", "Supposed Child Image")
	heritageCmd.Flags().StringSliceVarP(&Parents, "parent", "p", []string{}, "Supposed Parent Image/s. Can take multiple parents")
	heritageCmd.Flags().BoolVarP(&Relaxed, "relaxed", "r", false, "Do a relaxed check; only one parent has to match")
	heritageCmd.Flags().SortFlags = false
	heritageCmd.MarkFlagRequired("parent")
	heritageCmd.MarkFlagRequired("child")
}

func Strict(child string, parents ...string) (bool, error) {
	// Need to implement

	return false, nil
}
