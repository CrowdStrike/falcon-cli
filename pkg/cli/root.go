/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	shortDesc = `The CrowdStrike Falcon CLI`

	longDesc = templates.LongDesc(
		`The CrowdStrike Falcon CLI allows you to work effortlessly
        with the CrowdStrike Falcon platform.
    `)
)

// rootCmd represents the base command when called without any subcommands
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "falcon",
		Short: shortDesc,
		Long:  longDesc,
		RunE:  runHelp,
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.falcon/falcon.yaml)")
	cmd.Flags().Bool("version", false, "Show version")
	cmd.PersistentFlags().Bool("help", false, "Show help for command")

	return cmd
}

func runHelp(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
