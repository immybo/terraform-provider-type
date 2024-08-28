# Terraform Provider Type

This provider allows for stricter type checking than Terraform natively does.

In Terraform, you can specify a complex variable type, e.g.

```
variable "complex" {
  type = object({
    a = string
    b = optional(string)
  })
}
```

This allows you to specify values that must be specified on the type - Terraform will throw an error if you don't specify 'a' in the above type. However, there is no way to verify that you haven't specified values that are not specified on the type. If you were to specify 'c' in the above type, Terraform would simply ignore it. In situations where very complex variables are specified by contributors against a fast moving project this can cause issues if optional values are moved around or if values are removed.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

This provider currently provides a single data source which checks a JSON object against a JSONSchema to mitigate this. The intended use is to specify a non-constrained object as the Terraform variable:

```
variable "complex" {
  type = object
}
```

And then use this provider to ensure validity against a JSONSchema:

```
data "type_validation_json" {
  json_schema = <<EOF
{
  "type": "object",
  "required": ["a"],
  "properties": {
    "a": {
      "type": "string"
    },
    "b": {
      "type": "string"
    }
	}
}
EOF
	json_object = jsonencode(var.complex)
  fail_on_validation_error = true
}
```

The schema could also be loaded from a file.

If fail_on_validation_error is not provided, the data source will store any validation errors in validation_errors and set the is_valid flag to false instead of throwing an error.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
