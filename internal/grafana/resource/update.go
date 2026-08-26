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

func (r *AdaptiveTracesRawPolicy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
		ID:        plan.ID.ValueString(),
		Name:      plan.Name.ValueString(),
		Type:      plan.Type.ValueString(),
		Body:      json.RawMessage(plan.Body.ValueString()),
		ExpiresAt: expiresAt,
	}

	updatedPolicy, err := r.client.UpdatePolicy(ctx, policy)
	if err != nil {
		resp.Diagnostics.AddError("Failed updating policy", err.Error())
		return
	}

	plan.ID = types.StringValue(updatedPolicy.ID)
	plan.Type = types.StringValue(updatedPolicy.Type)

	// Canonicalize the response body before setting state to avoid unnecessary diffs
	if len(updatedPolicy.Body) > 0 {
		canonical, err := canonicalizeJSON(updatedPolicy.Body)
		if err != nil {
			resp.Diagnostics.AddError("Failed canonicalizing response body", err.Error())
			return
		}
		plan.Body = jsontypes.NewNormalizedValue(canonical)
	} else {
		plan.Body = jsontypes.NewNormalizedNull()
	}

	if updatedPolicy.ExpiresAt != "" {
		var diags diag.Diagnostics
		plan.ExpiresAt, diags = timetypes.NewRFC3339Value(updatedPolicy.ExpiresAt)
		resp.Diagnostics.Append(diags...)
	} else {
		plan.ExpiresAt = timetypes.NewRFC3339Null()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
