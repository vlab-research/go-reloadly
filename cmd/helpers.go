package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vlab-research/go-reloadly/reloadly"
)

func LoadTopupsService(cmd *cobra.Command) (*reloadly.Service, error) {
	svc := reloadly.NewTopups()
	return loadService(cmd, svc)
}

func LoadGiftCardsService(cmd *cobra.Command) (*reloadly.Service, error) {
	svc := reloadly.NewGiftCards()
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

func PrettyPrint(object interface{}) error {
	b, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
