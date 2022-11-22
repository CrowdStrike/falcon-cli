package config

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
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

	Selector string
}

func NewCmdConfig(f *utils.Factory) *cobra.Command {
	opts := &ConfigOptions{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "config <profile>",
		Short: "Configures a profile to use with CrowdStrike Falcon API",
		Long: templates.LongDesc(`
		Configure falcon CLI with CrowdStrike Falcon API.
		`),
		Aliases: []string{"login", "init"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			opts.Config = cfg

			if len(args) > 0 {
				opts.Selector = args[1]
			}

			if !opts.IO.CanPrompt() {
				return fmt.Errorf(heredoc.Doc(`
				 Interactive mode is disabled in this terminal.

				 Please run this command in an interactive terminal.
				`))
			}

			return configRun(opts)
		},
	}
	utils.DisableAuthCheck(cmd)

	return cmd
}

func configRun(opts *ConfigOptions) error {
	var qs = []*survey.Question{
		{
			Name: "clientId",
			Prompt: &survey.Password{
				Message: "Enter your CrowdStrike API Client ID",
			},
		},
		{
			Name: "clientSecret",
			Prompt: &survey.Password{
				Message: "Enter your CrowdStrike API Client Secret",
			},
		},
		{
			Name: "cid",
			Prompt: &survey.Input{
				Message: "Enter your CrowdStrike Customer ID (CID)",
			},
		},
		{
			Name: "memberCid",
			Prompt: &survey.Input{
				Message: "Enter your CrowdStrike Member CID",
			},
		},
		{
			// TODO: Should store valid options somewhere else perhaps use gofalcon
			Name: "cloud",
			Prompt: &survey.Select{
				Message: "Select your CrowdStrike Cloud",
				Options: []string{"us-1", "us-2", "eu-1"},
			},
		},
	}

	if opts.Selector == "" {
		// prompt for profile name
		qs = append(qs, &survey.Question{
			Name: "profile",
			Prompt: &survey.Input{
				Message: "What is the name of the profile you want to configure?",
				Default: "default",
			},
		})
	}

	err := survey.Ask(qs, &opts.Config)

	// TODO: Validate and write config

	return err
}
