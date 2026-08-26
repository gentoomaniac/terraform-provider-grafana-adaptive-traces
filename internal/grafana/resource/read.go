package resource

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r *AdaptiveTracesRawPolicy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tfPolicy Policy
	resp.Diagnostics.Append(req.State.Get(ctx, &tfPolicy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.client.GetPolicy(ctx, tfPolicy.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed reading policy", err.Error())
		return
	}

	data, err := json.Marshal(policy)
	if err != nil {
		resp.Diagnostics.AddError("failed marshalling data", err.Error())
		return
	}

	tflog.Debug(ctx, "testing read", map[string]any{"policy": string(data)})

	tfPolicy.Name = types.StringValue(policy.Name)
	tfPolicy.Type = types.StringValue(policy.Type)

	normalizedBody, err := canonicalizeJSON(policy.Body)
	if err != nil {
		resp.Diagnostics.AddError("failed normalizing body", err.Error())
		return
	}
	tfPolicy.Body = jsontypes.NewNormalizedValue(normalizedBody)

	if policy.ExpiresAt != "" {
		var diags diag.Diagnostics
		tfPolicy.ExpiresAt, diags = timetypes.NewRFC3339Value(policy.ExpiresAt)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		tfPolicy.ExpiresAt = timetypes.NewRFC3339Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &tfPolicy)...)
}

// canonicalizeJSON sorts keys alphabetically and removes all indentation/newlines,
// making the string byte-for-byte identical to HCL's jsonencode().
func canonicalizeJSON(raw []byte) (string, error) {
	var v any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber() // Preserves numeric types without float64 conversion
	if err := dec.Decode(&v); err != nil {
		return "", err
	}
	canonical, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(canonical), nil
}
