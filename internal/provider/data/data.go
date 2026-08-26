package data

import "github.com/gentoomaniac/terraform-provider-grafana-adaptive-traces/internal/grafana"

type ProviderData struct {
	Client *grafana.AdaptiveTracesClient
}
