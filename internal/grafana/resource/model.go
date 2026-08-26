package resource

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Policy struct {
	ID        types.String         `tfsdk:"id"`
	Name      types.String         `tfsdk:"name"`
	Type      types.String         `tfsdk:"type"`
	Body      jsontypes.Normalized `tfsdk:"body"`
	ExpiresAt timetypes.RFC3339    `tfsdk:"expires_at"`
}
