# reloadly

Work in progress.

## Usage

``` go
svc := reloadly.New()
svc.Auth(id, secret)

operator := svc.AutoDetect(mobile, country)
```

Also, for sandbox usage, you have:

``` go
svc := reloadly.NewSandbox()
```

Or you can create the `Service` type directly.

## Docs
