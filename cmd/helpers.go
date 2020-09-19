package cmd

import (
	"os"

	"github.com/nandanrao/go-reloadly/reloadly"
	"github.com/spf13/cobra"
)

func LoadService(cmd *cobra.Command) (*reloadly.Service, error) {
	sandbox, err := cmd.Flags().GetBool("sandbox")
	if err != nil {
		return nil, err
	}

	svc := reloadly.New()
	if sandbox {
		svc.Sandbox()
	}

	err = svc.Auth(os.Getenv("RELOADLY_ID"), os.Getenv("RELOADLY_SECRET"))
	if err != nil {
		return nil, err
	}
	return svc, nil
}
