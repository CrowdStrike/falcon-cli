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
package factory

import (
	"github.com/crowdstrike/falcon-cli/internal/config"
	"github.com/crowdstrike/falcon-cli/pkg/iostreams"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/crowdstrike/gofalcon/falcon/client"
)

func New(appVersion string) *utils.Factory {
	f := &utils.Factory{
		Config: configFunc(),
	}

	f.FalconClient = falconClientFunc(f, appVersion) // Depends on Config
	f.IOStreams = ioStreams(f)

	return f
}

func configFunc() func() (*config.Config, error) {
	return func() (*config.Config, error) {
		config, err := config.NewConfig()
		return &config, err
	}

}

func falconClientFunc(f *utils.Factory, appVersion string) func() (*client.CrowdStrikeAPISpecification, error) {
	return func() (*client.CrowdStrikeAPISpecification, error) {
		cfg, err := f.Config()

		if err != nil {
			return nil, err
		}

		client, err := falcon.NewClient(cfg.ApiConfig(appVersion))
		return client, err
	}
}

func ioStreams(f *utils.Factory) *iostreams.IOStreams {
	i := &iostreams.IOStreams{}
	io := i.NewIOStreams()

	// TODO: Implement a way to opt out of prompts (perhaps a FALCON_PROMPT_DISABLE env var)
	io.SetNeverPrompt(false)

	return io
}
