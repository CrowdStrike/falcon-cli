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
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/MakeNowJust/heredoc"
	"github.com/crowdstrike/falcon-cli/pkg/cmd/root"
	"github.com/crowdstrike/falcon-cli/pkg/config"
	"github.com/crowdstrike/falcon-cli/pkg/factory"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/crowdstrike/falcon-cli/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Run() error {
	cmdFactory := factory.New(version.Version)
	rootCmd := root.NewCmdRoot(cmdFactory, version.Version)

	cfg, err := cmdFactory.Config()

	stderr := cmdFactory.IOStreams.ErrOut
	if err != nil {
		fmt.Fprintf(stderr, "Error loading config: %v", err)
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

		formatter := &log.TextFormatter{}
		formatter.TimestampFormat = "2006-01-02 15:04:05"
		formatter.FullTimestamp = true
		formatter.DisableLevelTruncation = true
		formatter.DisableColors = false
		formatter.ForceColors = true
		log.SetFormatter(formatter)

		if viper.GetBool("verbose") {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug logging is set")
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

	config.ConfigFile = v.ConfigFileUsed()

	v.SetEnvPrefix("falcon")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return bindFlags(cmd, v)
}

// bindFlags binds the flags to the viper config
func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	err := v.BindPFlag("profile", cmd.Flags().Lookup("profile"))
	if err != nil {
		return err
	}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == "profile" {
			err := v.BindEnv(f.Name, "FALCON_PROFILE")

			if err != nil {
				return
			}
		}

		// change the flag name to snake_case
		viperKey := strings.ReplaceAll(f.Name, "-", "_")

		// add the profile flag in front of the viper key
		viperKey = fmt.Sprintf("%s.%s", v.GetString("profile"), viperKey)

		err = v.BindPFlag(viperKey, f)
		if err != nil {
			fmt.Printf("Error binding flag %s: %v", f.Name, err)
		}

		// bind env var over the new viper key (profile.flag)
		err = v.BindEnv(viperKey, fmt.Sprintf("FALCON_%s", strings.ToUpper(f.Name)))
		if err != nil {
			fmt.Printf("Error binding env var %s: %v", f.Name, err)
		}
	})

	return nil
}

func authHelp() string {
	return heredoc.Doc(`
		Authentication is required for this command. Please use 'falcon auth config' to configure your credentials.

		For more information, run: 'falcon auth config --help'
		`)
}
