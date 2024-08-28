// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTypeValidateJsonDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: validTypeValidateJson,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.type_validate_json.test",
						tfjsonpath.New("is_valid"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				Config: invalidTypeValidateJsonMissingRequiredAttributes,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.type_validate_json.test",
						tfjsonpath.New("is_valid"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"data.type_validate_json.test",
						tfjsonpath.New("validation_errors"),
						knownvalue.StringRegexp(regexp.MustCompile("missing properties 'title', 'director'")),
					),
				},
			},
			{
				Config:      invalidTypeValidateJsonMissingRequiredAttributesFailOnError,
				ExpectError: regexp.MustCompile("missing properties 'title', 'director'"),
			},
			{
				Config:      invalidTypeValidateJsonSchemaDoesNotCompile,
				ExpectError: regexp.MustCompile("at '/properties/title/type': value must be one of 'array', 'boolean', 'integer', 'null', 'number', 'object', 'string'"),
			},
			{
				Config:      invalidTypeValidateJsonObjectDoesNotCompile,
				ExpectError: regexp.MustCompile("Unable to read data, got error: invalid character 'u'"),
			},
		},
	})
}

const validTypeValidateJson = `
data "type_validate_json" "test" {
  json_schema = <<EOF
{
  "$id": "https://example.com/movie.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A representation of a movie",
  "type": "object",
  "required": ["title", "director"],
  "properties": {
    "title": {
      "type": "string"
    },
    "director": {
      "type": "string"
    }
	}
}
EOF
	json_object = <<EOF
{
	"title": "Sample Movie",
	"director": "John Director"
}
EOF
}
`

const invalidTypeValidateJsonMissingRequiredAttributes = `
data "type_validate_json" "test" {
  json_schema = <<EOF
{
  "$id": "https://example.com/movie.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A representation of a movie",
  "type": "object",
  "required": ["title", "director"],
  "properties": {
    "title": {
      "type": "string"
    },
    "director": {
      "type": "string"
    }
	}
}
EOF
	json_object = <<EOF
{
}
EOF
}
`

const invalidTypeValidateJsonMissingRequiredAttributesFailOnError = `
data "type_validate_json" "test" {
  json_schema = <<EOF
{
  "$id": "https://example.com/movie.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A representation of a movie",
  "type": "object",
  "required": ["title", "director"],
  "properties": {
    "title": {
      "type": "string"
    },
    "director": {
      "type": "string"
    }
	}
}
EOF
	json_object = <<EOF
{
}
EOF
	fail_on_validation_error = true
}
`

const invalidTypeValidateJsonSchemaDoesNotCompile = `
data "type_validate_json" "test" {
  json_schema = <<EOF
{
  "$id": "https://example.com/movie.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A representation of a movie",
  "type": "object",
  "required": ["title", "director"],
  "properties": {
    "title": {
      "type": "notarealtype"
    },
    "director": {
      "type": "string"
    }
	}
}
EOF
	json_object = <<EOF
{
	"title": "Sample Movie",
	"director": "John Director"
}
EOF
}
`

const invalidTypeValidateJsonObjectDoesNotCompile = `
data "type_validate_json" "test" {
  json_schema = <<EOF
{
  "$id": "https://example.com/movie.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A representation of a movie",
  "type": "object",
  "required": ["title", "director"],
  "properties": {
    "title": {
      "type": "string"
    },
    "director": {
      "type": "string"
    }
	}
}
EOF
	json_object = <<EOF
{
	"title": "Sample Movie",
	"director": "John Director",
	unquoted string
}
EOF
}
`
