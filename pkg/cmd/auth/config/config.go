package config

import (
	"fmt"

	"github.com/crowdstrike/falcon-cli/internal/config"
	"github.com/crowdstrike/falcon-cli/pkg/iostreams"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

type ConfigOptions struct {
	IO          *iostreams.IOStreams
	Config      config.Config
	Interactive bool
}

func NewCmdConfig(f *utils.Factory) *cobra.Command {
	opts := &ConfigOptions{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure falcon CLI with CrowdStrike Falcon API",
		Long: templates.LongDesc(`
		Configure falcon CLI with CrowdStrike Falcon API.
		`),
		Aliases: []string{"login", "init"},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			opts.Config = cfg

			if len(args) == 0 {
				opts.Interactive = true
			}

			if opts.Interactive && !opts.IO.CanPrompt() {
				return fmt.Errorf("client_id and client_secret must be provided as arguments when not running interactively")
			}

			return configRun(opts)
		},
	}

	return cmd
}

func configRun(opts *ConfigOptions) error {
	fmt.Println("configRun")
	fmt.Println(opts)
	fmt.Println(opts.Config.CID)
	fmt.Println(opts.Config.ClientID)
	fmt.Println(opts.Config.ClientSecret)
	fmt.Println(opts.Config.MemberCID)
	fmt.Println(opts.Config.Cloud)

	return nil
}
