package cmd

import (
	"os"

	"github.com/nandanrao/go-reloadly/reloadly"
	"github.com/spf13/cobra"
)

func LoadService(cmd *cobra.Command) (*reloadly.Service, error) {
	sb, err := cmd.Flags().GetBool("sandbox")
	if err != nil {
		return nil, err
	}
	var svc *reloadly.Service
	if sb {
		svc = reloadly.NewSandbox()
	} else {
		svc = reloadly.New()
	}

	err = svc.Auth(os.Getenv("RELOADLY_ID"), os.Getenv("RELOADLY_SECRET"))
	if err != nil {
		return nil, err
	}
	return svc, nil
}
