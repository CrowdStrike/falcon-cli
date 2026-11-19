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

package kubernetes

import (
	"github.com/crowdstrike/falcon-cli/pkg/cmdutil"
	"github.com/crowdstrike/falcon-cli/pkg/factory"
	"github.com/crowdstrike/falcon-cli/pkg/output"
	"github.com/crowdstrike/gofalcon/falcon/client/kubernetes_protection"
	"github.com/crowdstrike/gofalcon/falcon/models"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var containerTableDef = &output.TableDefinition{
	Headers: []string{"ID", "CLUSTER", "NAMESPACE", "IMAGE_ID", "CLOUD", "FIRST_SEEN"},
	RowFunc: func(item any) []string {
		r, ok := item.(*models.ModelsContainer)
		if !ok {
			return nil
		}
		return []string{strp(r.ContainerID), r.ClusterName, r.Namespace, r.ContainerImageID, strp(r.CloudName), r.FirstSeen}
	},
}

func NewCmdContainers(f *factory.Factory) *cobra.Command {
	opts := struct {
		filter string
		limit  int64
		offset int64
		out    string
		Table  cmdutil.TableFormattingOptions
		Client cmdutil.ClientFilterOptions
	}{limit: 25}

	cmd := &cobra.Command{
		Use:   "containers",
		Short: "List Kubernetes containers",
		Long:  templates.LongDesc(`List Kubernetes containers monitored by Falcon.`),
		Example: templates.Examples(`
            # List containers
            falcon fcs kubernetes containers

            # Filter by cluster
            falcon fcs kubernetes containers --filter "cluster_name:'my-cluster'"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			falconClient, err := f.FalconClient()
			if err != nil {
				return err
			}

			p := kubernetes_protection.NewContainerCombinedParams()
			p.Filter = cmdutil.StringPtr(opts.filter)
			p.Limit = &opts.limit
			if opts.offset > 0 {
				p.Offset = &opts.offset
			}

			resp, err := falconClient.KubernetesProtection.ContainerCombined(p)
			if err != nil {
				return cmdutil.HandleAPIError(err, "ContainerCombined")
			}

			resources := output.FilterAndSort(resp.Payload.Resources, containerTableDef, opts.Client.ToFilterOptions())
			tableOpts := opts.Table.ToOutputTableOptions()
			printer := output.NewPrinter(output.Format(opts.out), containerTableDef, tableOpts)
			return printer.Print(f.IOStreams.Out, resources)
		},
	}

	cmdutil.AddFilterFlag(cmd, &opts.filter)
	cmdutil.AddPaginationFlags(cmd, &opts.limit, &opts.offset, 25)
	cmdutil.AddOutputFlag(cmd, &opts.out)
	cmdutil.AddTableFormattingFlags(cmd, &opts.Table)
	cmdutil.AddClientFilterFlags(cmd, &opts.Client)
	return cmd
}
