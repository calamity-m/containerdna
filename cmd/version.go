package cmd

import (
	"fmt"

	"github.com/calamity-m/containerdna/pkg/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display version info",
		Long:  `Displays build and version info for this binary, including dependencies.`,
		Run: func(cmd *cobra.Command, args []string) {
			v := version.GetVersion()
			fmt.Printf("Module Version: %s\n", v.ModuleVersion)
			fmt.Printf("Built using Go Version: %s\n", v.GoVersion)
			fmt.Printf("Built using Git Commit: %s\n", v.GitCommit)
			fmt.Printf("Built using Commit Time: %s\n", v.CommitTime)
			for _, kv := range v.Dependencies {
				fmt.Printf("Dependency: %s - version: %s, hash: %s\n", kv.Path, kv.Version, kv.Sum)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
