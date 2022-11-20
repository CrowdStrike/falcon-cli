// Copyright (c) 2022 CrowdStrike, Inc.
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

package cli

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/crowdstrike/falcon-cli/internal/build"
	"github.com/crowdstrike/falcon-cli/internal/config"
	"github.com/crowdstrike/falcon-cli/pkg/cmd/factory"
	"github.com/crowdstrike/falcon-cli/pkg/cmd/root"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run() error {
	cmdFactory := factory.New(build.Version)
	rootCmd := root.NewCmdRoot(cmdFactory, build.Version)
	cobra.OnInitialize(initConfig)

	cfg, err := cmdFactory.Config()

	stderr := cmdFactory.IOStreams.ErrOut
	if err != nil {
		fmt.Fprintf(stderr, "Error loading config: %v", err)
	}

	// Support falcon help <command>
	if len(os.Args) > 1 && os.Args[1] == "help" {
		if len(os.Args) > 2 {
			os.Args[1] = os.Args[2]
		} else {
			os.Args = append(os.Args, "--help")
		}
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Do auth check if the command requires authentication
		if utils.IsAuthEnabled(cmd) && !utils.CheckAuth(*cfg) {
			return fmt.Errorf(authHelp())
		}
		return nil
	}

	return rootCmd.Execute()
}

func initConfig() {
	cfgFile := viper.GetString("config")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		config.ConfigFile = cfgFile
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		falconHome := fmt.Sprintf("%s/.falcon", home)
		// Search config in home directory with name "falcon" (without extension).
		viper.AddConfigPath(falconHome)
		viper.SetConfigType("yaml")
		viper.SetConfigName("falcon")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("falcon")
}

func authHelp() string {
	return heredoc.Doc(`
		Authentication is required for this command. Please use 'falcon auth config' to configure your credentials.

		For more information, run: 'falcon auth config --help'
		`)
}
