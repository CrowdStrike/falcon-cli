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

package init

import (
	log "github.com/sirupsen/logrus"

	"github.com/crowdstrike/falcon-cli/pkg/api"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/crowdstrike/gofalcon/pkg/falcon_util"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	shortDescInit = `Initialize the Falcon CLI tool`
	longDescInit  = templates.LongDesc(`Initialize the Falcon CLI tool`)
	examplesInit  = templates.Examples(`
	    # Initialize the CrowdStrike Falcon CLI tool
	    falcon init
		
		# Initialize the CrowdStrike Falcon CLI tool with defining OAuth2 client ID and secret
		falcon init --client-id <client_id> --client-secret <client_secret>
    `)

	cid, clientID, clientSecret, memberCID, cloud string
)

// NewInitCmd represents the init command
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   shortDescInit,
		Long:    longDescInit,
		Example: examplesInit,
		RunE:    runInit,
	}
	cmd.Flags().StringVarP(&cid, "cid", "f", "", "The Falcon Customer ID (CID)")
	cmd.Flags().StringVarP(&clientID, "client-id", "u", "", "The Falcon API Oauth client ID")
	cmd.Flags().StringVarP(&clientSecret, "client-secret", "s", "", "The Falcon API Oauth client secret")
	cmd.Flags().StringVarP(&memberCID, "member-cid", "m", "", "The Falcon API member CID")
	cmd.Flags().StringVarP(&cloud, "cloud", "r", "autodiscover", "The Falcon API Cloud Region")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	if cid != "" {
		if !utils.ValidateRegExp("^[0-9a-fA-F]{32}-[0-9a-fA-F]{2}$", cid) {
			log.Errorf("Invalid member CID provided: %s", cid)
		}
	}

	if clientID == "" {
		clientID = falcon_util.PromptUser(`Please provide your CrowdStrike Falcon OAuth2 API Client ID`)
	}

	if !utils.ValidateRegExp("[0-9a-z]{32}", clientID) {
		log.Errorf("Invalid client ID provided: %s", clientID)
	}

	if clientSecret == "" {
		clientSecret = falcon_util.PromptUser(`Please provide your CrowdStrike Falcon OAuth2 API Client Secret`)
	}

	if !utils.ValidateRegExp("[0-9a-zA-Z]{40}", clientSecret) {
		log.Errorf("Invalid client ID provided: %s", clientSecret)
	}

	if memberCID != "" {
		if !utils.ValidateRegExp("^[0-9a-fA-F]{32}-[0-9a-fA-F]{2}$", memberCID) {
			log.Errorf("Invalid member CID provided: %s", memberCID)
		}
	}

	if _, err := falcon.CloudValidate(cloud); err != nil {
		log.Error(err)
	}

	config := api.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		MemberCID:    memberCID,
		Cloud:        cloud,
	}

	utils.ConfigExists(utils.ConfigFile)
	utils.WriteYAML(config, utils.ConfigFile)

	return nil
}
