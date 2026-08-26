package resource

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/gentoomaniac/terraform-provider-grafana-adaptive-traces/internal/grafana"
)

func (r *AdaptiveTracesRawPolicy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Policy

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var expiresAt string
	if !plan.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
		expiresAt = plan.ExpiresAt.ValueString()
	}
	policy := grafana.Policy{
		Name:      plan.Name.ValueString(),
		Type:      plan.Type.ValueString(),
		Body:      json.RawMessage(plan.Body.ValueString()),
		ExpiresAt: expiresAt,
	}

	createdPolicy, err := r.client.CreatePolicy(ctx, policy)
	if err != nil {
		resp.Diagnostics.AddError("Failed creating policy", err.Error())
		return
	}

	plan.ID = types.StringValue(createdPolicy.ID)
	plan.Type = types.StringValue(createdPolicy.Type)

	// Canonicalize the response body before setting state to avoid unnecessary diffs
	if len(createdPolicy.Body) > 0 {
		canonical, err := canonicalizeJSON(createdPolicy.Body)
		if err != nil {
			resp.Diagnostics.AddError("Failed canonicalizing response body", err.Error())
			return
		}
		plan.Body = jsontypes.NewNormalizedValue(canonical)
	} else {
		plan.Body = jsontypes.NewNormalizedNull()
	}

	if createdPolicy.ExpiresAt != "" {
		var diags diag.Diagnostics
		plan.ExpiresAt, diags = timetypes.NewRFC3339Value(createdPolicy.ExpiresAt)
		resp.Diagnostics.Append(diags...)
	} else {
		plan.ExpiresAt = timetypes.NewRFC3339Null()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
