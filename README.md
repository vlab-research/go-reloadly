# reloadly

A Go client library and CLI for the Reloadly API, supporting airtime topups and gift cards.

## Installation

### From GitHub

Install the CLI directly from GitHub:

```bash
# Install the latest version
go install github.com/vlab-research/go-reloadly/cmd@latest

# Or install a specific version
go install github.com/vlab-research/go-reloadly/cmd@v0.0.22
```

After installation, you can use the CLI directly:
```bash
reloadly-cli operators list IN
```

### From Source

Clone and build from source:

```bash
# Clone the repository
git clone https://github.com/vlab-research/go-reloadly.git
cd go-reloadly

# Build the CLI
go build -o reloadly-cli ./cmd

# Make it available system-wide (optional)
sudo cp reloadly-cli /usr/local/bin/
```

## CLI Usage

The reloadly CLI provides easy access to Reloadly's airtime topup services. Here are the most common operations:

### Quick Start

Set up your credentials (you can also use environment variables):
```bash
export RELOADLY_CLIENT_ID="your-client-id"
export RELOADLY_CLIENT_SECRET="your-client-secret"
```

### List Operators by Country

View all available operators for a specific country:

```bash
# List all operators in India
reloadly-cli operators list IN

# List all operators in Nigeria  
reloadly-cli operators list NG

# List all operators in Kenya
reloadly-cli operators list KE
```

This will show you:
- Operator ID and name
- Country and denomination type
- Bundle and data support
- Available fixed amounts (both USD and local currency)

### Get Operator Details

Get detailed information about a specific operator:

```bash
# Get details by operator name and country
reloadly-cli operators get IN "Airtel India"

# Get details by operator ID
reloadly-cli operators info 200
```

This shows comprehensive operator information including:
- Supported amounts and currencies
- Commission and discount rates
- FX rates and limits
- Available fixed amounts and descriptions

### Make Individual Topups

Perform airtime topups to mobile numbers:

```bash
# Simple topup with auto-detection
reloadly-cli topups single "+919876543210" 10.0 IN

# Topup with specific operator
reloadly-cli topups single "+919876543210" 10.0 IN --operator "Airtel India"

# Topup with tolerance for amount matching
reloadly-cli topups single "+919876543210" 10.0 IN --tolerance 2.0
```

### Check Amount Support

Verify which operators support a specific amount:

```bash
# Check which operators support 10 USD with 2 USD tolerance
reloadly-cli operators amounts IN 10.0 2.0
```

### Sandbox Mode

For testing, use sandbox mode:

```bash
export RELOADLY_SANDBOX=true
reloadly-cli operators list IN
```

## Library Usage

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

