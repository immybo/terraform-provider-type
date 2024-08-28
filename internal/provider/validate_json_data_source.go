// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TypeValidateJsonDataSource{}

func NewTypeValidateJsonDataSource() datasource.DataSource {
	return &TypeValidateJsonDataSource{}
}

type TypeValidateJsonDataSource struct {
}

type TypeValidateJsonDataSourceModel struct {
	JsonSchema            types.String `tfsdk:"json_schema"`
	JsonObject            types.String `tfsdk:"json_object"`
	IsValid               types.Bool   `tfsdk:"is_valid"`
	ValidationErrors      types.String `tfsdk:"validation_errors"`
	FailOnValidationError types.Bool   `tfsdk:"fail_on_validation_error"`
}

func (d *TypeValidateJsonDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_validate_json"
}

func (d *TypeValidateJsonDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source that can be used to validate the type of a Terraform object against a YAML schema.",

		Attributes: map[string]schema.Attribute{
			"json_schema": schema.StringAttribute{
				MarkdownDescription: "The expected JSON schema of the provided object - must be valid JSONSchema.",
				Required:            true,
			},
			"json_object": schema.StringAttribute{
				MarkdownDescription: "The object to check the type of - must be valid JSON.",
				Required:            true,
			},
			"is_valid": schema.BoolAttribute{
				MarkdownDescription: "Whether or not the provided object is valid according to the provided schema.",
				Computed:            true,
			},
			"validation_errors": schema.StringAttribute{
				MarkdownDescription: "A human-readable string containing one or more validation errors. This will always be empty if is_valid is true.",
				Computed:            true,
			},
			"fail_on_validation_error": schema.BoolAttribute{
				MarkdownDescription: "Whether or not a Terraform error should be thrown if any validation errors are discovered.",
				Optional:            true,
			},
		},
	}
}

func (d *TypeValidateJsonDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TypeValidateJsonDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schema, err := jsonschema.UnmarshalJSON(strings.NewReader(data.JsonSchema.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Schema Error", fmt.Sprintf("Unable to read schema, got error: %s", err))
		return
	}

	jsonData, err := jsonschema.UnmarshalJSON(strings.NewReader(data.JsonObject.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Schema Error", fmt.Sprintf("Unable to read data, got error: %s", err))
		return
	}

	schemaRef := "schema.json"

	compiler := jsonschema.NewCompiler()
	err = compiler.AddResource(schemaRef, schema)

	if err != nil {
		resp.Diagnostics.AddError("Schema Error", fmt.Sprintf("Unable to add JSON schema to compiler, got error: %s", err))
		return
	}

	compiledSchema, err := compiler.Compile(schemaRef)
	if err != nil {
		resp.Diagnostics.AddError("Schema Error", fmt.Sprintf("Unable to compile schema, got error: %s", err))
		return
	}

	validationErr := compiledSchema.Validate(jsonData)

	if data.FailOnValidationError.ValueBool() && validationErr != nil {
		resp.Diagnostics.AddError("Validation Error", validationErr.Error())
	}

	if validationErr != nil {
		data.ValidationErrors = types.StringValue(validationErr.Error())
	}

	data.IsValid = types.BoolValue(validationErr == nil)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
