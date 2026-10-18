// Copyright (c) 2026 CrowdStrike, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package config

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/crowdstrike/falcon-cli/pkg/config"
	"github.com/crowdstrike/falcon-cli/pkg/factory"
	"github.com/crowdstrike/falcon-cli/pkg/iostreams"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var falconClouds = []string{"autodiscover", "us-1", "us-2", "eu-1", "us-gov-1"}

type Options struct {
	IO          *iostreams.IOStreams
	Config      config.Config
	Interactive bool

	Selector string
}

func NewCmdConfig(f *factory.Factory) *cobra.Command {
	opts := &Options{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "config [profile]",
		Short: "Configures a profile to use with CrowdStrike Falcon API",
		Long: templates.LongDesc(`
		Configure falcon CLI with CrowdStrike Falcon API.
		`),
		Aliases: []string{"login", "init"},
		Args:    cobra.MaximumNArgs(1),
		Example: templates.Examples(`
			# Configure the default profile interactively
			falcon auth config

			# Configure a named profile non-interactively
			falcon auth config prod --client-id <ID> --client-secret <SECRET>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Selector = args[0]
			} else if cmd.Flags().Changed("profile") {
				opts.Selector, _ = cmd.Flags().GetString("profile")
			}

			// Non-interactive: credentials supplied via flags
			if cmd.Flags().Changed("client-id") && cmd.Flags().Changed("client-secret") {
				opts.Config = configFromFlags(cmd)
				if opts.Selector == "" {
					opts.Selector = "default"
				}
				return saveConfig(opts)
			}

			if !opts.IO.CanPrompt() {
				return errors.New(heredoc.Doc(`
				 Interactive mode is disabled in this terminal.

				 Please run this command in an interactive terminal, or pass
				 --client-id and --client-secret.
				`))
			}

			if err := configRun(opts); err != nil {
				return err
			}
			return saveConfig(opts)
		},
	}
	utils.DisableAuthCheck(cmd)

	return cmd
}

func configFromFlags(cmd *cobra.Command) config.Config {
	get := func(name string) string {
		v, _ := cmd.Flags().GetString(name)
		return v
	}
	return config.Config{
		CID:          get("cid"),
		ClientID:     get("client-id"),
		ClientSecret: get("client-secret"),
		MemberCID:    get("member-cid"),
		Cloud:        get("cloud"),
		BaseURL:      get("base-url"),
	}
}

func saveConfig(opts *Options) error {
	if opts.Selector != "" {
		opts.Config.Profile = opts.Selector
	}
	if opts.Config.ClientID == "" || opts.Config.ClientSecret == "" {
		return errors.New("client ID and client secret are required")
	}

	path := config.ConfigFile
	if path == "" {
		var err error
		if path, err = config.DefaultConfigFile(); err != nil {
			return err
		}
	}

	if err := opts.Config.Save(path); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Fprintf(opts.IO.ErrOut, "Saved profile %q to %s\n", opts.Config.Profile, path)
	return nil
}

func configRun(opts *Options) error {
	var qs = []*survey.Question{
		{
			Name: "clientId",
			Prompt: &survey.Password{
				Message: "Enter your CrowdStrike API Client ID:",
			},
		},
		{
			Name: "clientSecret",
			Prompt: &survey.Password{
				Message: "Enter your CrowdStrike API Client Secret:",
			},
		},
		{
			Name: "cid",
			Prompt: &survey.Input{
				Message: "Enter your CrowdStrike Customer ID (CID):",
			},
		},
		{
			Name: "memberCid",
			Prompt: &survey.Input{
				Message: "Enter your CrowdStrike Member CID:",
			},
		},
		{
			// TODO: Should store valid options somewhere else perhaps use gofalcon
			Name: "cloud",
			Prompt: &survey.Select{
				Message: "Select your CrowdStrike Cloud:",
				Options: falconClouds,
				Default: "autodiscover",
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

	return survey.Ask(qs, &opts.Config)
}
