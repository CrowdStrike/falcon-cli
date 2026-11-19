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

package registration

import (
	"fmt"
	"strings"

	"github.com/crowdstrike/falcon-cli/pkg/cmdutil"
	"github.com/crowdstrike/falcon-cli/pkg/factory"
	"github.com/crowdstrike/falcon-cli/pkg/output"
	"github.com/crowdstrike/gofalcon/falcon/client/cloud_google_cloud_registration"
	"github.com/crowdstrike/gofalcon/falcon/client/cspm_registration"
	"github.com/crowdstrike/gofalcon/falcon/models"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewGCPCmd represents the fcs registration gcp command group.
func NewGCPCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gcp",
		Short: "Manage GCP account registrations",
		Long:  templates.LongDesc(`List, register, and delete GCP cloud account registrations.`),
	}

	cmd.AddCommand(
		newCmdGCPList(f),
		newCmdGCPCreate(f),
		newCmdGCPDelete(f),
	)
	return cmd
}

var gcpAccountTableDef = &output.TableDefinition{
	Headers: []string{"DISPLAY_NAME", "PROJECT_ID", "CSPM", "STATUS"},
	RowFunc: func(item any) []string {
		r, ok := item.(*models.DomainGCPAccountV1)
		if !ok || r == nil {
			return nil
		}
		cspm := "disabled"
		if r.CspmEnabled != nil && *r.CspmEnabled {
			cspm = "enabled"
		}
		status := "registered"
		if len(r.Conditions) > 0 && r.Conditions[0] != nil {
			if r.Conditions[0].Status != nil {
				status = *r.Conditions[0].Status
			}
		}
		return []string{
			r.DisplayName,
			r.ProjectID,
			cspm,
			status,
		}
	},
}

func newCmdGCPList(f *factory.Factory) *cobra.Command {
	var out string
	var limit int64

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List registered GCP projects",
		Long:  templates.LongDesc(`List all GCP projects registered with Falcon Cloud Security.`),
		RunE: func(_ *cobra.Command, _ []string) error {
			client, err := f.FalconClient()
			if err != nil {
				return err
			}

			params := cspm_registration.NewGetCSPMGCPAccountParams()
			params.Limit = &limit

			// Migrating to CloudRegistrationGcpGetEntities changes the output
			// shape; deferred to a separate change.
			resp, _, err := client.CspmRegistration.GetCSPMGCPAccount(params) //nolint:staticcheck // SA1019: migration deferred
			if err != nil {
				return cmdutil.HandleAPIError(err, "list GCP registrations")
			}

			printer := output.NewPrinter(output.Format(out), gcpAccountTableDef)
			if resp == nil || resp.Payload == nil {
				return printer.Print(f.IOStreams.Out, nil)
			}
			return printer.Print(f.IOStreams.Out, resp.Payload.Resources)
		},
	}

	cmd.Flags().Int64Var(&limit, "limit", 100, "Maximum number of results")
	cmdutil.AddOutputFlag(cmd, &out)
	return cmd
}

func newCmdGCPCreate(f *factory.Factory) *cobra.Command {
	opts := struct {
		entityID         string
		infraProjectID   string
		deploymentMethod string
	}{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Register a GCP project",
		Long:  templates.LongDesc(`Register a new GCP project/organization with Falcon Cloud Security.`),
		Example: templates.Examples(`
            # Register a GCP project
            falcon fcs registration gcp create --entity-id <project-number> --infra-project-id <project-id> --deployment-method service_account
        `),
		RunE: func(_ *cobra.Command, _ []string) error {
			client, err := f.FalconClient()
			if err != nil {
				return err
			}

			params := cloud_google_cloud_registration.NewCloudRegistrationGcpCreateRegistrationParams()
			params.Body = &models.DtoGCPRegistrationCreateRequestExtV1{
				Resources: []*models.DtoCreateGCPRegistrationRequest{
					{
						EntityID:         []string{opts.entityID},
						InfraProjectID:   &opts.infraProjectID,
						DeploymentMethod: &opts.deploymentMethod,
					},
				},
			}

			resp, err := client.CloudGoogleCloudRegistration.CloudRegistrationGcpCreateRegistration(params)
			if err != nil {
				return cmdutil.HandleAPIError(err, "create GCP registration")
			}

			if resp != nil && resp.Payload != nil && len(resp.Payload.Resources) > 0 {
				fmt.Fprintf(f.IOStreams.Out, "Registered GCP project: %s\n", resp.Payload.Resources[0].InfraProjectID)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.entityID, "entity-id", "", "GCP project number or organization ID (required)")
	cmd.Flags().StringVar(&opts.infraProjectID, "infra-project-id", "", "Infrastructure project ID (required)")
	cmd.Flags().StringVar(&opts.deploymentMethod, "deployment-method", "service_account", "Deployment method (service_account, workload_identity)")
	_ = cmd.MarkFlagRequired("entity-id")
	_ = cmd.MarkFlagRequired("infra-project-id")
	return cmd
}

func newCmdGCPDelete(f *factory.Factory) *cobra.Command {
	var ids []string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete GCP registrations",
		Long:  templates.LongDesc(`Delete one or more GCP registrations from Falcon Cloud Security.`),
		RunE: func(_ *cobra.Command, _ []string) error {
			client, err := f.FalconClient()
			if err != nil {
				return err
			}

			// The GCP delete API takes a single comma-separated string of IDs
			params := cloud_google_cloud_registration.NewCloudRegistrationGcpDeleteRegistrationParams()
			params.Ids = strings.Join(ids, ",")

			_, err = client.CloudGoogleCloudRegistration.CloudRegistrationGcpDeleteRegistration(params)
			if err != nil {
				return cmdutil.HandleAPIError(err, "delete GCP registration")
			}

			fmt.Fprintf(f.IOStreams.Out, "Deleted %d GCP registration(s)\n", len(ids))
			return nil
		},
	}

	cmdutil.AddIDsFlag(cmd, &ids, true)
	return cmd
}
