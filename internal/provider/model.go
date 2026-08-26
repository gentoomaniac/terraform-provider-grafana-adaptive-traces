package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Model struct {
	URL    types.String `tfsdk:"url"`
	APIKey types.String `tfsdk:"api_key"`
}
