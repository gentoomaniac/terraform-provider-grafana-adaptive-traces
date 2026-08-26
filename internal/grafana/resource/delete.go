package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *AdaptiveTracesRawPolicy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Policy
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	err := r.client.DeletePolicy(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed deleting policy", err.Error())
		return
	}
}
