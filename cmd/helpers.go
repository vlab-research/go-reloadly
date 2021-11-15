package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vlab-research/go-reloadly/reloadly"
)

func LoadTopupsService(cmd *cobra.Command) (*reloadly.Service, error) {
	svc := reloadly.NewTopups()
	return loadService(cmd, svc)
}


func loadService(cmd *cobra.Command, svc *reloadly.Service) (*reloadly.Service, error) {
	sandbox, err := cmd.Flags().GetBool("sandbox")
	if err != nil {
		return nil, err
	}

	if sandbox {
		svc.Sandbox()
	}

	err = svc.Auth(os.Getenv("RELOADLY_ID"), os.Getenv("RELOADLY_SECRET"))
	if err != nil {
		return nil, err
	}

	fmt.Println("Authorized with Reloadly")

	return svc, nil
}
