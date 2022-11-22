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

	config "github.com/crowdstrike/falcon-cli/pkg/cmd/init"
	"github.com/crowdstrike/falcon-cli/pkg/cmd/sensor"
	"github.com/crowdstrike/falcon-cli/pkg/cmd/version"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	ver "github.com/crowdstrike/falcon-cli/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	commands = []*cobra.Command{
		config.NewInitCmd(),
		version.VersionCmd(),
		sensor.SensorCmd(),
	}
	cfgFile  string
	replacer = strings.NewReplacer("-", "_", ".", "_")
)

type CLI struct {
	// Root command name.
	commandName string

	// Root command.
	cmd *cobra.Command
}

func Run() error {
	c, _ := CreateCLIAndRoot()
	return c.cmd.Execute()
}

func CreateCLIAndRoot() (*CLI, *cobra.Command) {

	c := &CLI{}
	c.cmd = newRootCmd()
	c.commandName = "falcon"

	cobra.OnInitialize(initConfig)

	// Add the subcommands
	err := c.addSubCommands()
	if err != nil {
		log.Fatal(err)
	}

	root := commands[0].Root()

	root.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if err := viper.BindPFlag(strings.ReplaceAll(flag.Name, "-", "_"), flag); err != nil {
			return
		}
	})

	c.cmd.PersistentPreRun = rootPersistentPreRun

	return c, root
}

// addSubCommands adds the additional commands.
func (c *CLI) addSubCommands() error {
	for _, cmd := range commands {
		for _, subCmd := range c.cmd.Commands() {
			if cmd.Name() == subCmd.Name() {
				return fmt.Errorf("command %q already exists", cmd.Name())
			}
		}
		c.cmd.AddCommand(cmd)
	}
	return nil
}

func rootPersistentPreRun(cmd *cobra.Command, args []string) {
	if ok, err := cmd.Flags().GetBool("version"); err == nil && ok {
		fmt.Println(ver.VersionString())
		os.Exit(0)
	}

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		utils.ConfigFile = cfgFile
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

	viper.SetEnvPrefix("falcon")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}
