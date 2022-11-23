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

package sensor

import (
	downloadCmd "github.com/crowdstrike/falcon-cli/pkg/cmd/sensor/download"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	shortDesc = `Manage the CrowdStrike Falcon Sensor`
	longDesc  = templates.LongDesc(`Manage the CrowdStrike Falcon Sensor`)
	examples  = templates.Examples(`
        # Download the CrowdStrike Falcon Sensor
        falcon sensor download
    `)
)

// NewCmdSensor represents the sensor command
func NewSensorCmd(f *utils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sensor",
		Short:   shortDesc,
		Long:    longDesc,
		Example: examples,
	}

	cmd.AddCommand(
		downloadCmd.NewCmdDownload(f),
	)
	return cmd
}
