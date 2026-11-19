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

	"github.com/crowdstrike/falcon-cli/pkg/cmdutil"
	"github.com/crowdstrike/falcon-cli/pkg/factory"
	"github.com/crowdstrike/falcon-cli/pkg/output"
	"github.com/crowdstrike/gofalcon/falcon/client/cloud_azure_registration"
	"github.com/crowdstrike/gofalcon/falcon/client/cspm_registration"
	"github.com/crowdstrike/gofalcon/falcon/models"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewAzureCmd represents the fcs registration azure command group.
func NewAzureCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "azure",
		Short: "Manage Azure account registrations",
		Long:  templates.LongDesc(`List, register, update, and delete Azure cloud account registrations.`),
	}

	cmd.AddCommand(
		newCmdAzureList(f),
		newCmdAzureCreate(f),
		newCmdAzureDelete(f),
	)
	return cmd
}

var azureAccountTableDef = &output.TableDefinition{
	Headers: []string{"TENANT_ID", "ACCOUNT_TYPE", "CLIENT_ID", "STATUS"},
	RowFunc: func(item any) []string {
		r, ok := item.(*models.RegistrationAzureAccountV1Ext)
		if !ok || r == nil {
			return nil
		}
		status := "registered"
		if len(r.Conditions) > 0 && r.Conditions[0] != nil {
			if r.Conditions[0].Status != nil {
				status = *r.Conditions[0].Status
			}
		}
		tenantID := ""
		if r.TenantID != "" {
			tenantID = r.TenantID
		}
		return []string{
			tenantID,
			r.AccountType,
			r.ClientID,
			status,
		}
	},
}

func newCmdAzureList(f *factory.Factory) *cobra.Command {
	var out string
	var limit int64

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List registered Azure tenants",
		Long:  templates.LongDesc(`List all Azure tenants registered with Falcon Cloud Security.`),
		RunE: func(_ *cobra.Command, _ []string) error {
			client, err := f.FalconClient()
			if err != nil {
				return err
			}

			params := cspm_registration.NewGetCSPMAzureAccountParams()
			params.Limit = &limit

			// The replacement (cloud-security-registration-azure GET) requires a
			// tenant or registration ID and cannot list all tenants.
			resp, _, err := client.CspmRegistration.GetCSPMAzureAccount(params) //nolint:staticcheck // SA1019: no list-all replacement yet
			if err != nil {
				return cmdutil.HandleAPIError(err, "list Azure registrations")
			}

			printer := output.NewPrinter(output.Format(out), azureAccountTableDef)
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

func newCmdAzureCreate(f *factory.Factory) *cobra.Command {
	opts := struct {
		tenantID string
	}{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Register an Azure tenant",
		Long:  templates.LongDesc(`Register a new Azure tenant with Falcon Cloud Security.`),
		Example: templates.Examples(`
            # Register an Azure tenant
            falcon fcs registration azure create --tenant-id <uuid>
        `),
		RunE: func(_ *cobra.Command, _ []string) error {
			client, err := f.FalconClient()
			if err != nil {
				return err
			}

			params := cloud_azure_registration.NewCloudRegistrationAzureCreateRegistrationParams()
			params.Body = &models.AzureAzureRegistrationCreateRequestExtV1{
				Resource: &models.AzureAzureRegistrationCreateInput{
					TenantID: &opts.tenantID,
				},
			}

			resp, err := client.CloudAzureRegistration.CloudRegistrationAzureCreateRegistration(params)
			if err != nil {
				return cmdutil.HandleAPIError(err, "create Azure registration")
			}

			if resp != nil && resp.Payload != nil && len(resp.Payload.Resources) > 0 {
				r := resp.Payload.Resources[0]
				tenantID := ""
				if r.TenantID != nil {
					tenantID = *r.TenantID
				}
				fmt.Fprintf(f.IOStreams.Out, "Registered Azure tenant: %s\n", tenantID)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.tenantID, "tenant-id", "", "Azure tenant ID (required)")
	_ = cmd.MarkFlagRequired("tenant-id")
	return cmd
}

func newCmdAzureDelete(f *factory.Factory) *cobra.Command {
	var tenantIDs []string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete Azure tenant registrations",
		Long:  templates.LongDesc(`Delete one or more Azure tenant registrations from Falcon Cloud Security.`),
		Example: templates.Examples(`
            # Delete an Azure registration
            falcon fcs registration azure delete --tenant-ids <uuid>
        `),
		RunE: func(_ *cobra.Command, _ []string) error {
			client, err := f.FalconClient()
			if err != nil {
				return err
			}

			params := cloud_azure_registration.NewCloudRegistrationAzureDeleteRegistrationParams()
			params.TenantIds = tenantIDs

			_, err = client.CloudAzureRegistration.CloudRegistrationAzureDeleteRegistration(params)
			if err != nil {
				return cmdutil.HandleAPIError(err, "delete Azure registration")
			}

			fmt.Fprintf(f.IOStreams.Out, "Deleted %d Azure registration(s)\n", len(tenantIDs))
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&tenantIDs, "tenant-ids", nil, "Azure tenant IDs to delete (required)")
	_ = cmd.MarkFlagRequired("tenant-ids")
	return cmd
}
