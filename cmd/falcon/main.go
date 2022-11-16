package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/crowdstrike/falcon-cli/pkg/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
