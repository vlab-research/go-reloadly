package reloadly

import (
	"fmt"
	"net/http"
)

type ProductsPage struct {
	Content []Product `json:"content,omitempty"`
	Page    int64     `json:"page,omitempty"`
	Size    int64     `json:"size,omitempty"`
}

type Product struct {
	ProductID                   int64     `json:"productId,omitempty"`
	ProductName                 string    `json:"productName,omitempty"`
	Global                      bool      `json:"global,omitempty"`
	SenderFee                   float64   `json:"senderFee,omitempty"`
	DiscountPercentage          float64   `json:"discountPercentage,omitempty"`
	DenominationType            string    `json:"denominationType,omitempty"`
	RecipientCurrencyCode       string    `json:"recipientCurrencyCode,omitempty"`
	MinRecipientDenomination    float64   `json:"minRecipientDenomination,omitempty"`
	MaxRecipientDenomination    float64   `json:"maxRecipientDenomination,omitempty"`
	SenderCurrencyCode          string    `json:"senderCurrencyCode,omitempty"`
	MinSenderDenomination       float64   `json:"minSenderDenomination,omitempty"`
	MaxSenderDenomination       float64   `json:"maxSenderDenomination,omitempty"`
	FixedRecipientDenominations []float64 `json:"fixedRecipientDenominations,omitempty"`
	FixedSenderDenominations    []float64 `json:"fixedSenderDenominations,omitempty"`
	FixedRecipientToSender      int64     `json:"fixedRecipientToSender,omitempty"`
	LogoUrls                    []string  `json:"logoUrls,omitempty"`
	BrandID                     int64     `json:"brandId,omitempty"`
	BrandName                   string    `json:"brandName,omitempty"`
	IsoName                     string    `json:"isoName,omitempty"`
	Name                        string    `json:"name,omitempty"`
	CountryCode                 string    `json:"countryCode,omitempty"`
}

type RedeemInstructions struct {
	BrandID   int64  `json:"brandId,omitempty"`
	BrandName string `json:"brandName,omitempty"`
	Concise   string `json:"concise,omitempty"`
	Verbose   string `json:"verbose,omitempty"`
}

type DiscountsPage struct {
	Content []Discount `json:"content,omitempty"`
	Page    int64      `json:"page,omitempty"`
	Size    int64      `json:"size,omitempty"`
}

type Discount struct {
	*Product           `json:"product,omitempty"`
	DiscountPercentage float64 `json:"discountPercentage,omitempty"`
}

type TransactionsPage struct {
	Content []Transaction `json:"content,omitempty"`
	Page    int64         `json:"page,omitempty"`
	Size    int64         `json:"size,omitempty"`
}

type Transaction struct {
	TransactionId          int64   `json:"transactionId,omitempty"`
	Amount                 float64 `json:"amount,omitempty"`
	Discount               float64 `json:"discount,omitempty"`
	CurrencyCode           string  `json:"currencyCode,omitempty"`
	Fee                    float64 `json:"fee,omitempty"`
	RecipientEmail         string  `json:"recipientEmail,omitempty"`
	CustomIdentifier       string  `json:"customIdentifier,omitempty"`
	Status                 string  `json:"status,omitempty"`
	TransactionCreatedTime string  `json:"transactionCreatedTime,omitempty"`
	ProductID              int64   `json:"productId,omitempty"`
	ProductName            string  `json:"productName,omitempty"`
	CountryCode            string  `json:"countryCode,omitempty"`
	Quantity               int64   `json:"quantity,omitempty"`
	UnitPrice              float64 `json:"unitPrice,omitempty"`
	TotalPrice             float64 `json:"totalPrice,omitempty"`
	BrandID                int64   `json:"brandId,omitempty"`
	BrandName              string  `json:"brandName,omitempty"`
}

type GiftCardOrder struct {
	ProductID        int64   `json:"productId,omitempty"`
	CountryCode      string  `json:"countryCode,omitempty"`
	Quantity         int64   `json:"quantity,omitempty"`
	UnitPrice        float64 `json:"unitPrice,omitempty"`
	CustomIdentifier string  `json:"customIdentifier,omitempty"`
	SenderName       string  `json:"senderName,omitempty"`
	RecipientEmail   string  `json:"recipientEmail,omitempty"`
}

type Card struct {
	CardNumber string `json:"cardNumber,omitempty"`
	PinCode    string `json:"pinCode,omitempty"`
}

type GiftCardsService struct {
	*Service
	acceptHeader string
}

func NewGiftCards() *Service {
	return &Service{
		http.DefaultClient,
		"https://giftcards.reloadly.com",
		"https://auth.reloadly.com",
		nil,
		"",
		"",
		"https://giftcards-sandbox.reloadly.com",
	}
}

func (s *Service) GiftCards() *GiftCardsService {
	return &GiftCardsService{s, "application/com.reloadly.giftcards-v1+json"}
}

func (s *GiftCardsService) Products(page int64, size int64) (ProductsPage, error) {
	path := fmt.Sprintf("/products?page=%v&size=%v", page, size)
	resp := new(ProductsPage)
	_, err := s.Request("GET", path, nil, resp, s.acceptHeader)
	resp.Page = page
	return *resp, err
}

func (s *GiftCardsService) Product(productId int64) (Product, error) {
	path := fmt.Sprintf("/products/%v", productId)
	resp := new(Product)
	_, err := s.Request("GET", path, nil, resp, s.acceptHeader)
	return *resp, err
}

func (s *GiftCardsService) ProductsByCountry(country string) ([]Product, error) {
	path := fmt.Sprintf("/countries/%v/products", country)
	resp := new([]Product)
	_, err := s.Request("GET", path, nil, resp, s.acceptHeader)
	return *resp, err
}

func (s *GiftCardsService) RedeemInstructions() ([]RedeemInstructions, error) {
	resp := new([]RedeemInstructions)
	_, err := s.Request("GET", "/redeem-instructions", nil, resp, s.acceptHeader)
	return *resp, err
}

func (s *GiftCardsService) RedeemInstructionsByBrand(brandId int64) (RedeemInstructions, error) {
	path := fmt.Sprintf("/redeem-instructions/%v", brandId)
	resp := new(RedeemInstructions)
	_, err := s.Request("GET", path, nil, resp, s.acceptHeader)
	return *resp, err
}

func (s *GiftCardsService) Discounts(page int64, size int64) (DiscountsPage, error) {
	path := fmt.Sprintf("/discounts?page=%v&size=%v", page, size)
	resp := new(DiscountsPage)
	_, err := s.Request("GET", path, nil, resp, s.acceptHeader)
	resp.Page = page
	return *resp, err
}

func (s *GiftCardsService) DiscountByProduct(productId int64) (Discount, error) {
	path := fmt.Sprintf("/products/%v/discounts", productId)
	resp := new(Discount)
	_, err := s.Request("GET", path, nil, resp, s.acceptHeader)
	return *resp, err
}

func (s *GiftCardsService) Transactions(page int64, size int64) (TransactionsPage, error) {
	path := fmt.Sprintf("/reports/transactions?page=%v&size=%v", page, size)
	resp := new(TransactionsPage)
	_, err := s.Request("GET", path, nil, resp, s.acceptHeader)
	resp.Page = page
	return *resp, err
}

func (s *GiftCardsService) Transaction(transactionId int64) (Transaction, error) {
	path := fmt.Sprintf("/reports/transactions/%v", transactionId)
	resp := new(Transaction)
	_, err := s.Request("GET", path, nil, resp, s.acceptHeader)
	return *resp, err
}

func (s *GiftCardsService) Order(order GiftCardOrder) (Transaction, error) {
	resp := new(Transaction)
	_, err := s.Request("POST", "/orders", order, resp, s.acceptHeader)
	return *resp, err
}

func (s *GiftCardsService) GetRedeemCode(transactionId int64) ([]Card, error) {
	path := fmt.Sprintf("/orders/transactions/%v/cards", transactionId)
	resp := new([]Card)
	_, err := s.Request("GET", path, nil, resp, s.acceptHeader)
	return *resp, err
}
