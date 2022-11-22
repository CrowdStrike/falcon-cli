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
	"log"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/crowdstrike/falcon-cli/internal/build"
	"github.com/crowdstrike/falcon-cli/internal/config"
	"github.com/crowdstrike/falcon-cli/pkg/cmd/factory"
	"github.com/crowdstrike/falcon-cli/pkg/cmd/root"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Run() error {
	cmdFactory := factory.New(build.Version)
	rootCmd := root.NewCmdRoot(cmdFactory, build.Version)

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
		err := initConfig(cmd)

		if err != nil {
			return err
		}

		if err = viper.GetViper().BindPFlags(cmd.PersistentFlags()); err != nil {
			log.Fatalf("Error binding flags to viper: %v", err)
		}
		//Do auth check if the command requires authentication
		if utils.IsAuthCheckEnabled(cmd) && !utils.CheckAuth(cfg) {
			return fmt.Errorf(authHelp())
		}
		return nil
	}

	return rootCmd.Execute()
}

func initConfig(cmd *cobra.Command) error {
	v := viper.GetViper()

	cfgFile := v.GetString("config")

	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
		config.ConfigFile = cfgFile
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		falconHome := fmt.Sprintf("%s/.falcon", home)
		// Search config in home directory with name "falcon" (without extension).
		v.AddConfigPath(falconHome)
		v.SetConfigType("yaml")
		v.SetConfigName("config")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("Error reading config file: %v", err)
		}
	}

	v.SetEnvPrefix("falcon")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		viperKey := strings.ReplaceAll(f.Name, "-", "_")
		v.BindPFlag(viperKey, f)
	})
}

func authHelp() string {
	return heredoc.Doc(`
		Authentication is required for this command. Please use 'falcon auth config' to configure your credentials.

		For more information, run: 'falcon auth config --help'
		`)
}
