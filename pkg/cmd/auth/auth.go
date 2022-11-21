package auth

import (
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	authConfigCmd "github.com/crowdstrike/falcon-cli/pkg/cmd/auth/config"
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

func NewAuthCmd(f *utils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate falcon CLI with CrowdStrike Falcon API",
	}

	// Add subcommands
	utils.DisableAuthCheck(cmd)

	cmd.AddCommand(authConfigCmd.NewCmdConfig(f))

	return cmd
}
