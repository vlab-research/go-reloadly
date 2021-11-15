package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vlab-research/go-reloadly/reloadly"
)

var giftCardsCmd = &cobra.Command{
	Use:   "gift-cards",
	Short: "Send gift cards to mobile numbers",
	Long:  "Send gift cards to mobile numbers",
}

var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Get information on over 200 gift cards",
	Long:  "Get information on over 200 gift cards",
	RunE: func(cmd *cobra.Command, args []string) error {
		page, err := cmd.Flags().GetInt64("page")
		if err != nil {
			return err
		}

		size, err := cmd.Flags().GetInt64("size")
		if err != nil {
			return err
		}

		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		ps, err := svc.GiftCards().Products(page, size)
		if err != nil {
			return err
		}
		PrettyPrint(ps)

		return nil
	},
}

var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Access the details of a particular gift card",
	Long:  "Access the details of a particular gift card",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires product id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		productId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		p, err := svc.GiftCards().Product(productId)
		if err != nil {
			return err
		}
		PrettyPrint(p)

		return nil
	},
}

var productsByCountryCmd = &cobra.Command{
	Use:   "products-by-country",
	Short: "Retrieve details of every gift card that is available in a country",
	Long:  "Retrieve details of every gift card that is available in a country",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires country name")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		country := args[0]

		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		ps, err := svc.GiftCards().ProductsByCountry(country)
		if err != nil {
			return err
		}
		PrettyPrint(ps)

		return nil
	},
}

var redeemInstructionsCmd = &cobra.Command{
	Use:   "redeem-instructions",
	Short: "Provides details on how to redeem a gift card",
	Long:  "Provides details on how to redeem a gift card",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		ins, err := svc.GiftCards().RedeemInstructions()
		if err != nil {
			return err
		}
		PrettyPrint(ins)

		return nil
	},
}

var redeemInstructionsByBrandCmd = &cobra.Command{
	Use:   "redeem-instructions-by-brand",
	Short: "Retrieve the redeem instructions for a particular gift card brand",
	Long:  "Retrieve the redeem instructions for a particular gift card brand",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires brand id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		brandId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		ins, err := svc.GiftCards().RedeemInstructionsByBrand(brandId)
		if err != nil {
			return err
		}
		PrettyPrint(ins)

		return nil
	},
}

var discountsCmd = &cobra.Command{
	Use:   "discounts",
	Short: "Fetch data of every gift card that has an available discount at the point of purchase",
	Long:  "Fetch data of every gift card that has an available discount at the point of purchase",
	RunE: func(cmd *cobra.Command, args []string) error {
		page, err := cmd.Flags().GetInt64("page")
		if err != nil {
			return err
		}

		size, err := cmd.Flags().GetInt64("size")
		if err != nil {
			return err
		}

		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		ds, err := svc.GiftCards().Discounts(page, size)
		if err != nil {
			return err
		}
		PrettyPrint(ds)

		return nil
	},
}

var discountByProductCmd = &cobra.Command{
	Use:   "discount-by-product",
	Short: "Get the details of an active discount on a particular gift card",
	Long:  "Get the details of an active discount on a particular gift card",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires product id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		productId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		d, err := svc.GiftCards().DiscountByProduct(productId)
		if err != nil {
			return err
		}
		PrettyPrint(d)

		return nil
	},
}

var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Fetch details of every gift card purchase",
	Long:  "Fetch details of every gift card purchase",
	RunE: func(cmd *cobra.Command, args []string) error {
		page, err := cmd.Flags().GetInt64("page")
		if err != nil {
			return err
		}

		size, err := cmd.Flags().GetInt64("size")
		if err != nil {
			return err
		}

		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		ts, err := svc.GiftCards().Transactions(page, size)
		if err != nil {
			return err
		}
		PrettyPrint(ts)

		return nil
	},
}

var transactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Fetch the details of a particular gift card purchase",
	Long:  "Fetch the details of a particular gift card purchase",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires transaction id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		transactionId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		t, err := svc.GiftCards().Transaction(transactionId)
		if err != nil {
			return err
		}
		PrettyPrint(t)

		return nil
	},
}

var orderGiftCardCmd = &cobra.Command{
	Use:   "order",
	Short: "Purchase a gift card",
	Long:  "Purchase a gift card",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 7 {
			return errors.New("requires productId, countryCode, quantity, unitPrice, customIdentifier, senderName and recipientEmail")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		productId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		countryCode := args[1]

		quantity, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return err
		}

		unitPrice, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			return err
		}

		customIdentifier := args[4]
		senderName := args[5]
		recipientEmail := args[6]

		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		order := reloadly.GiftCardOrder{productId, countryCode, quantity, unitPrice, customIdentifier, senderName, recipientEmail}
		o, err := svc.GiftCards().Order(order)
		if err != nil {
			return err
		}
		PrettyPrint(o)

		return nil
	},
}

var getRedeemCodeCmd = &cobra.Command{
	Use:   "get-redeem-code",
	Short: "Retrieve details of an already purchased gift card",
	Long:  "Retrieve details of an already purchased gift card",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires transaction id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		transactionId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		svc, err := LoadGiftCardService(cmd)
		if err != nil {
			return err
		}

		r, err := svc.GiftCards().GetRedeemCode(transactionId)
		if err != nil {
			return err
		}
		PrettyPrint(r)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(giftCardsCmd)
	giftCardsCmd.AddCommand(productsCmd)
	productsCmd.Flags().Int64P("page", "", 1, "page number")
	productsCmd.Flags().Int64P("size", "", 10, "max number of results to return")
	productsCmd.SetArgs([]string{
		"--page=1",
		"--size=10",
	})
	giftCardsCmd.AddCommand(productCmd)
	giftCardsCmd.AddCommand(productsByCountryCmd)
	giftCardsCmd.AddCommand(redeemInstructionsCmd)
	giftCardsCmd.AddCommand(redeemInstructionsByBrandCmd)
	giftCardsCmd.AddCommand(discountsCmd)
	discountsCmd.Flags().Int64P("page", "", 1, "page number")
	discountsCmd.Flags().Int64P("size", "", 10, "max number of results to return")
	discountsCmd.SetArgs([]string{
		"--page=1",
		"--size=10",
	})
	giftCardsCmd.AddCommand(discountByProductCmd)
	giftCardsCmd.AddCommand(transactionsCmd)
	transactionsCmd.Flags().Int64P("page", "", 1, "page number")
	transactionsCmd.Flags().Int64P("size", "", 10, "max number of results to return")
	transactionsCmd.SetArgs([]string{
		"--page=1",
		"--size=10",
	})
	giftCardsCmd.AddCommand(transactionCmd)
	giftCardsCmd.AddCommand(orderGiftCardCmd)
	giftCardsCmd.AddCommand(getRedeemCodeCmd)
}
