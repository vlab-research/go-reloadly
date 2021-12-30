package reloadly

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/products.json")
	products := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/products", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, products)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.GiftCards().Products(1, 100)

	assert.Nil(t, err)
	assert.Equal(t, int64(5904), res.Content[0].ProductID)
	assert.Equal(t, "Microsoft 365 Personal", res.Content[0].ProductName)
}

func TestGetProductsPagination(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/products", r.URL.Path)
		assert.Equal(t, "page=2&size=1", r.URL.RawQuery)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.GiftCards().Products(2, 1)

	assert.NotNil(t, err)
}

func TestGetProduct(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/apple.json")
	product := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/products/10", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, product)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.GiftCards().Product(10)

	assert.Nil(t, err)
	assert.Equal(t, int64(10), res.ProductID)
	assert.Equal(t, "App Store & iTunes Austria", res.ProductName)
}

func TestGetProductReturnsError(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/products/100", r.URL.Path)

		w.WriteHeader(400)
		fmt.Fprintf(w, `{"timeStamp":"2021-11-16 20:28:00","message":"The product was either not found or is no longer available, Please contact support","path":"/products/100","errorCode":null,"infoLink":null,"details":[]}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.GiftCards().Product(100)

	assert.NotNil(t, err)
	assert.Equal(t, 400, err.(APIError).StatusCode)
}

func TestProductsByCountry(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/products_by_country.json")
	products := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/countries/es/products", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, products)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.GiftCards().ProductsByCountry("es")

	assert.Nil(t, err)
	assert.Equal(t, int64(11), res[1].ProductID)
	assert.Equal(t, "App Store & iTunes Spain", res[1].ProductName)
}

func TestProductsByCountryReturnsError(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/countries/ABC/products", r.URL.Path)

		w.WriteHeader(400)
		fmt.Fprintf(w, `{"timeStamp":"2021-11-16 20:25:17","message":"No products were found for the given country code. For a list of valid country codes visit https://www.nationsonline.org/oneworld/country_code_list.htm","path":"/countries/ABC/products","errorCode":null,"infoLink":null,"details":[]}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.GiftCards().ProductsByCountry("ABC")

	assert.NotNil(t, err)
	assert.Equal(t, 400, err.(APIError).StatusCode)
}

func TestRedeemInstructions(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/redeem_instructions.json")
	instructions := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/redeem-instructions", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, instructions)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.GiftCards().RedeemInstructions()

	assert.Nil(t, err)
	assert.Equal(t, int64(2), res[1].BrandID)
	assert.Equal(t, "Amazon", res[1].BrandName)
}

func TestRedeemInstructionsByBrand(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/apple_instructions.json")
	instructions := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/redeem-instructions/4", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, instructions)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.GiftCards().RedeemInstructionsByBrand(4)

	assert.Nil(t, err)
	assert.Equal(t, int64(4), res.BrandID)
	assert.Equal(t, "Apple Music", res.BrandName)
}

func TestRedeemInstructionsByBrandReturnsError(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/redeem-instructions/66", r.URL.Path)

		w.WriteHeader(400)
		fmt.Fprintf(w, `{"timeStamp":"2021-11-16 20:21:55","message":"The redeem instruction was either not found or is no longer available, Please contact support","path":"/redeem-instructions/66","errorCode":null,"infoLink":null,"details":[]}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.GiftCards().RedeemInstructionsByBrand(66)

	assert.NotNil(t, err)
	assert.Equal(t, 400, err.(APIError).StatusCode)
}

func TestDiscounts(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/discounts.json")
	discounts := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/discounts", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, discounts)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.GiftCards().Discounts(1, 100)

	assert.Nil(t, err)
	assert.Equal(t, "1-800-PetSupplies", res.Content[0].Product.ProductName)
	assert.Equal(t, int64(2), res.Content[1].Product.ProductID)
}

func TestDiscountsPagination(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/discounts", r.URL.Path)
		assert.Equal(t, "page=2&size=1", r.URL.RawQuery)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.GiftCards().Discounts(2, 1)

	assert.NotNil(t, err)
}

func TestDiscountByProduct(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/products/2/discounts", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, `{"product":{"productId":2,"productName":"Amazon UK","countryCode":"GB","global":false},"discountPercentage":0.7}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.GiftCards().DiscountByProduct(2)

	assert.Nil(t, err)
	assert.Equal(t, "Amazon UK", res.Product.ProductName)
	assert.Equal(t, 0.7, res.DiscountPercentage)
}

func TestDiscountByProductReturnsNotFoundError(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/products/25/discounts", r.URL.Path)

		w.WriteHeader(404)
		fmt.Fprintf(w, `{"timeStamp":"2021-11-16 17:28:15","message":"Commission resource not found","path":"/products/2000/discounts","errorCode":null,"infoLink":null,"details":[]}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.GiftCards().DiscountByProduct(25)

	assert.NotNil(t, err)
	assert.Equal(t, 404, err.(APIError).StatusCode)
}

func TestTransactions(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/transactions.json")
	transactions := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reports/transactions", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, transactions)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.GiftCards().Transactions(1, 100)

	assert.Nil(t, err)
	assert.Equal(t, "INR", res.Content[0].CurrencyCode)
	assert.Equal(t, "yeahyeah", res.Content[0].CustomIdentifier)
}

func TestTransactionsPagination(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reports/transactions", r.URL.Path)
		assert.Equal(t, "page=2&size=1", r.URL.RawQuery)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.GiftCards().Transactions(2, 1)

	assert.NotNil(t, err)
}

func TestTransaction(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/transaction.json")
	transaction := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reports/transactions/563", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, transaction)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.GiftCards().Transaction(563)

	assert.Nil(t, err)
	assert.Equal(t, "INR", res.CurrencyCode)
	assert.Equal(t, 37.16255, res.Fee)
}

func TestTransactionNotFoundError(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reports/transactions/1000", r.URL.Path)

		w.WriteHeader(404)
		fmt.Fprint(w, `{"timeStamp":"2021-11-16 07:37:17","message":"Gift Card transaction not found","path":"/reports/transactions/1000","errorCode":null,"infoLink":null,"details":[]}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.GiftCards().Transaction(1000)

	assert.NotNil(t, err)
	assert.Equal(t, 404, err.(APIError).StatusCode)
}

func TestOrder(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/transaction.json")
	transaction := string(dat)

	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/orders", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, transaction)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	order := GiftCardOrder{157, "US", 1, 59.99, "test-card", "John Doe", "test@test.com", "test-id"}
	res, err := svc.GiftCards().Order(order)

	assert.Nil(t, err)
	assert.Equal(t, int64(563), res.TransactionId)
	assert.Equal(t, "test@test.com", res.RecipientEmail)
}

func TestOrderRequiresCertainFieldsAndNotOthers(t *testing.T) {
	// does not require custom identifier
	j := `{"productId": 157, "countryCode": "US", "quantity": 1, "UnitPrice": 5, "senderName": "John Doe", "recipientEmail": "test@test.com"}`

	order := new(GiftCardOrder)
	json.Unmarshal([]byte(j), &order)

	validate := validator.New()
	err := validate.Struct(order)
	assert.Nil(t, err)

	// requires productId
	j = `{"countryCode": "US", "quantity": 1, "UnitPrice": 5, "senderName": "John Doe", "recipientEmail": "test@test.com"}`
	order = new(GiftCardOrder)
	json.Unmarshal([]byte(j), &order)
	err = validate.Struct(order)
	assert.NotNil(t, err)
}

func TestOrderReturnsError(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/orders/transactions/100/cards", r.URL.Path)

		w.WriteHeader(400)
		fmt.Fprint(w, `{"timeStamp":"2021-11-16 07:31:37","message":"Insufficient funds in the wallet to complete this transaction","path":"/orders","errorCode":"INSUFFICIENT_BALANCE","infoLink":null,"details":[]}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.GiftCards().GetRedeemCode(100)

	assert.NotNil(t, err)
	assert.Equal(t, 400, err.(APIError).StatusCode)
	assert.Equal(t, "INSUFFICIENT_BALANCE", err.(APIError).ErrorCode)
}

func TestGetRedeemCode(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/orders/transactions/100/cards", r.URL.Path)

		w.Header().Set("Content-Type", "application/com.reloadly.giftcards-v1+json")
		fmt.Fprintf(w, `[{"cardNumber": "ABC-XYZ"}]`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	res, err := svc.GiftCards().GetRedeemCode(100)

	assert.Nil(t, err)
	assert.Equal(t, "ABC-XYZ", res[0].CardNumber)
}

func TestGetRedeemCodeReturnsNotFoundError(t *testing.T) {
	ts, _ := TestServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/orders/transactions/100/cards", r.URL.Path)

		w.WriteHeader(404)
		fmt.Fprint(w, `{"timeStamp":"2021-11-16 07:24:29","message":"Invalid transaction id","path":"/orders/transactions/100/cards","errorCode":null,"infoLink":null,"details":[]}`)
	})

	svc := &Service{BaseUrl: ts.URL, Client: &http.Client{}}
	_, err := svc.GiftCards().GetRedeemCode(100)

	assert.NotNil(t, err)
	assert.Equal(t, 404, err.(APIError).StatusCode)
}
