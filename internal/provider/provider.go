package provider

import (
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/gentoomaniac/terraform-provider-grafana-adaptive-traces/internal/grafana"
	raw_resource "github.com/gentoomaniac/terraform-provider-grafana-adaptive-traces/internal/grafana/resource"
	"github.com/gentoomaniac/terraform-provider-grafana-adaptive-traces/internal/provider/data"
)

// Ensure AdaptiveTracesProvider satisfies various provider interfaces.
var _ provider.Provider = &AdaptiveTracesProvider{}

type AdaptiveTracesProvider struct{}

func New() provider.Provider {
	return &AdaptiveTracesProvider{}
}

func (p *AdaptiveTracesProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "grafana"
}

func (p *AdaptiveTracesProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A provider for Grafana Adaptive Traces",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional:    true,
				Description: "Grafana Cloud's API URL. May alternatively be set via the GRAFANA_AT_API_URL environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Tenant ID and Access Policy Token (or API key) for Grafana Cloud in the format ':'. May alternatively be set via the GRAFANA_AT_API_KEY environment variable.",
			},
		},
	}
}

func (p *AdaptiveTracesProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config Model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.URL.IsNull() || config.URL.IsUnknown() {
		config.URL = types.StringValue(os.Getenv("GRAFANA_AT_API_URL"))
	}
	if config.APIKey.IsNull() || config.APIKey.IsUnknown() {
		config.APIKey = types.StringValue(os.Getenv("GRAFANA_AT_API_KEY"))
	}

	if config.URL.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Tenant Url not set",
			"A Tenant Url must be provided. To find your URL value, go to your Grafana Cloud account and check the Details page of your hosted Tempo endpoint.",
		)
		return
	}

	rawAuth := config.APIKey.ValueString()
	authParts := strings.SplitN(rawAuth, ":", 2)
	if len(authParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid API Key Format",
			"The api_key must be provided in the format 'tenantId:token'. The token needs adaptive-traces:admin scope for your stack ID.",
		)
		return
	}

	tenantId := authParts[0]
	apiKey := authParts[1]

	parsedURL, err := url.Parse(config.URL.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid URL", err.Error())
		return
	}

	client, err := grafana.NewAdaptiveTracesClient(parsedURL, tenantId, apiKey)
	if err != nil {
		resp.Diagnostics.AddError("Client Initialization Failed", err.Error())
		return
	}

	pd := &data.ProviderData{
		Client: client,
	}

	resp.ResourceData = pd
	resp.DataSourceData = pd
}

func (p *AdaptiveTracesProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *AdaptiveTracesProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{raw_resource.NewAdaptiveTracesRawPolicy}
}
