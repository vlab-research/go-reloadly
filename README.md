# reloadly

Work in progress.

## Usage

### Getting Started

``` go
svc := reloadly.New()
svc.Auth(id, secret)

operator := svc.OperatorsAutoDetect("+34987467293", "ES")
```

Also, for sandbox usage, you have:

``` go
svc := reloadly.New()
svc.Sandbox()

// now all calls use the sandbox base url
svc.Auth(id, secret)
operator := svc.OperatorsAutoDetect("+34987467293", "ES")
```

Or you can create the `Service` type directly.

### Making requests

``` go
req := &TopupRequest{
    RecipientPhone: &RecipientPhone{operator.Country.IsoName, "+9187654467"},
    OperatorID: operator.OperatorID,
	Amount: amount,
}

respV := new(TopupResponse)
httpResponse, err := svc.Request("POST", "/topups", req, respV)

// respV will have the marshalled json response
// err will be an APIError, if we get a json error response, or
// an error created directly from the http library or marshaling
```

### Convenience Methods

Some convenience methods for simple tasks:

``` go
svc.Topup("+3441983489", operator, 10)
```

### Topups using suggested amounts

This is useful...
