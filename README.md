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

## Development

### Requirements

Before we get started, make sure you have Docker installed in your local computer. If you don't already have Docker installed, visit https://docs.docker.com/get-docker.

### Local environment

This project supports hot reloading. Hot reloading allows you to run tests automatically while editing files, enabling an interactive TDD workflow.

To start the hot reload workflow, navigate to your terminal and type `./test.sh`. Once it is running, any editing of a file will trigger the test suite.

### Pre-commit hooks

To keep coding style consistent throughout the code base, this project uses a pre-commit hook, which runs style-correction libraries on every commit to the code base.

To setup the pre-commit hook, navigate to your terminal and type, `git config core.hooksPath $(pwd)`.

