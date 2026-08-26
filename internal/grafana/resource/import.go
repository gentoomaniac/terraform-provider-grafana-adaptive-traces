package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *AdaptiveTracesRawPolicy) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	policies, err := r.client.GetPolicies(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed fetching policies", err.Error())
		return
	}

	id := ""
	for _, policy := range policies {
		if policy.Name == req.ID {
			id = policy.ID
		}
	}

	if id == "" {
		resp.Diagnostics.AddError("No policy with name found: ", req.ID)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
