// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	fwflex "github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	ERNameParameter = "Ephemeral Resource Parameter"
)

// @EphemeralResource(aws_ssm_parameter, name="Parameter")
func newEphemeralParameter(_ context.Context) (ephemeral.EphemeralResourceWithConfigure, error) {
	return &ephemeralParameter{}, nil
}

type ephemeralParameter struct {
	framework.EphemeralResourceWithConfigure
}

func (e *ephemeralParameter) Metadata(_ context.Context, _ ephemeral.MetadataRequest, response *ephemeral.MetadataResponse) {
	response.TypeName = "aws_ssm_parameter"
}

func (e *ephemeralParameter) Schema(ctx context.Context, _ ephemeral.SchemaRequest, response *ephemeral.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			names.AttrARN: framework.ARNAttributeComputedOnly(),
			names.AttrName: schema.StringAttribute{
				Required: true,
			},
			names.AttrType: schema.StringAttribute{
				Computed: true,
			},
			names.AttrValue: schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			names.AttrVersion: schema.Int64Attribute{
				Computed: true,
			},
			"with_decryption": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (e *ephemeralParameter) Open(ctx context.Context, request ephemeral.OpenRequest, response *ephemeral.OpenResponse) {
	var data epParameterData
	conn := e.Meta().SSMClient(ctx)

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	input := ssm.GetParameterInput{}
	response.Diagnostics.Append(fwflex.Expand(ctx, data, &input)...)
	if response.Diagnostics.HasError() {
		return
	}

	output, err := findParameterByName(ctx, conn, *input.Name, *input.WithDecryption)
	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.SSM, create.ErrActionReading, ERNameParameter, data.ARN.String(), err),
			err.Error(),
		)
		return
	}

	data.Value = fwflex.StringValueToFramework(ctx, string(*output.Value))

	response.Diagnostics.Append(response.Result.Set(ctx, &data)...)
}

type epParameterData struct {
	ARN            types.String `tfsdk:"arn"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	Value          types.String `tfsdk:"value"`
	Version        types.Int64  `tfsdk:"version"`
	WithDecryption types.Bool   `tfsdk:"with_decryption"`
}
