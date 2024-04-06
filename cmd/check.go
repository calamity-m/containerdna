package cmd

import (
	"github.com/calamity-m/paternity/pkg/paternity"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	Parent string

	Child string

	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().Msg("Entering check command")
			if Parent == "" || Child == "" {
				return
			}
			paternity.Paternity(Parent, Child)
			//fmt.Println(cmd.Flags().GetString("parent"))
		},
	}
)

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&Parent, "parent", "p", "", "Supposed Parent Image")
	checkCmd.Flags().StringVarP(&Child, "child", "c", "", "Supposed Child Image")
	checkCmd.MarkFlagRequired("parent")
	checkCmd.MarkFlagRequired("child")
}
