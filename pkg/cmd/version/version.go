/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package version

import (
	"fmt"

	"github.com/crowdstrike/falcon-cli/pkg/version"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	shortDesc = `Print the falcon version`
	longDesc  = templates.LongDesc(`Print the falcon version`)
	examples  = templates.Examples(`
        # Print the Falcon CLI and GO version information
        falcon version
    `)
)

// versionCmd represents the version command
func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   shortDesc,
		Long:    longDesc,
		Example: examples,
		RunE:    runVer,
	}
	return cmd
}

func runVer(_ *cobra.Command, _ []string) error {
	fmt.Println(version.VersionString())
	return nil
}
