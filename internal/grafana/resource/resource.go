package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/gentoomaniac/terraform-provider-grafana-adaptive-traces/internal/grafana"
	provider_data "github.com/gentoomaniac/terraform-provider-grafana-adaptive-traces/internal/provider/data"
)

var _ resource.Resource = &AdaptiveTracesRawPolicy{}

func NewAdaptiveTracesRawPolicy() resource.Resource {
	return &AdaptiveTracesRawPolicy{}
}

type AdaptiveTracesRawPolicy struct {
	client *grafana.AdaptiveTracesClient
}

func (r *AdaptiveTracesRawPolicy) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_adaptive_traces_policy"
}

func (r *AdaptiveTracesRawPolicy) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*provider_data.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderData, got: %T", req.ProviderData),
		)
		return
	}

	r.client = data.Client
}
