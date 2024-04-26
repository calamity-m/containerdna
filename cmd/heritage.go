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
		Long: "Validates that a child orignates from the provided parent(s).\n" +
			"This is done by verifying the parent(s) 1..n layers correspond to 1..n layers in the child.",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Debug("Entering heritageCmd")

			valid, err := heritage.ValidateHeritage(Relaxed, Child, Parents...)
			if err != nil {
				fmt.Println("Encountered errors while attempting to evaluate heritage. Errors are:")
				fmt.Println(err)
				os.Exit(1)
			} else {
				fmt.Printf("Validity: %t", valid)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(heritageCmd)

	heritageCmd.Flags().StringVarP(&Child, "child", "c", "", "Child image to validate")
	heritageCmd.Flags().StringSliceVarP(&Parents, "parent", "p", []string{}, "Parent(s) for child to match against.")
	heritageCmd.Flags().BoolVarP(&Relaxed, "relaxed", "r", false, "Check for at least one parent to be the base of provided child")
	heritageCmd.Flags().SortFlags = false
	heritageCmd.MarkFlagRequired("parent")
	heritageCmd.MarkFlagRequired("child")
}
