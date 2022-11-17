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

	log "github.com/sirupsen/logrus"

	"github.com/crowdstrike/falcon-cli/internal/flags"
	"github.com/crowdstrike/falcon-cli/pkg/cmd/version"
	ver "github.com/crowdstrike/falcon-cli/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	commands = []*cobra.Command{
		version.VersionCmd(),
	}
	cfgFile string
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

	root.PersistentFlags().Bool(flags.Verbose, false, "Enable verbose logging")
	if err := viper.BindPFlags(root.PersistentFlags()); err != nil {
		log.Fatalf("Failed to bind %s flags: %v", root.Name(), err)
	}
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
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		falconHome := fmt.Sprintf("%s/.falcon", home)
		if err := os.Mkdir(falconHome, 0700); err != nil {
			if !os.IsExist(err) {
				log.Fatal(err)
			}
		}

		// Search config in home directory with name "falcon" (without extension).
		viper.AddConfigPath(falconHome)
		viper.SetConfigType("yaml")
		viper.SetConfigName("falcon")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}
