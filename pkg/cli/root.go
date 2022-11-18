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
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	shortDesc = `The CrowdStrike Falcon CLI`

	longDesc = templates.LongDesc(
		`The CrowdStrike Falcon CLI allows you to work effortlessly
        with the CrowdStrike Falcon platform.
    `)
	cid, clientID, clientSecret, memberCID, cloud string
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
	cmd.PersistentFlags().StringVarP(&cid, "cid", "f", "", "The Falcon Customer ID (CID)")
	cmd.PersistentFlags().StringVarP(&clientID, "client-id", "u", "", "The Falcon API Oauth client ID")
	cmd.PersistentFlags().StringVarP(&clientSecret, "client-secret", "s", "", "The Falcon API Oauth client secret")
	cmd.PersistentFlags().StringVarP(&memberCID, "member-cid", "m", "", "The Falcon API member CID")
	cmd.PersistentFlags().StringVarP(&cloud, "cloud", "r", "autodiscover", "The Falcon API Cloud Region")

	return cmd
}

func runHelp(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
