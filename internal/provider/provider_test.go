package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the HashiCups client is properly configured.
	// It is also possible to use the HASHICUPS_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "grafana" {}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"grafana": providerserver.NewProtocol6WithError(New()),
	}
)

func TestAdaptiveTracesPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "grafana_adaptive_traces_policy" "example" {
  name       = "test_rule"
  type       = "string_attribute"
  body       = jsonencode({
      key = "testfield.some-non-existing-field",
      values = [
          "ever"
      ]
  })
  expires_at = "2099-12-31T23:59:59Z"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("grafana_adaptive_traces_policy.example", "items.#", "1"),
					// Verify first order item
					resource.TestCheckResourceAttr("grafana_adaptive_traces_policy.example", "items.0.name", "test_rule"),
					resource.TestCheckResourceAttr("grafana_adaptive_traces_policy.example", "items.0.type", "string_attribute"),
					resource.TestCheckResourceAttr("grafana_adaptive_traces_policy.example", "items.0.body", "1"),
					resource.TestCheckResourceAttr("grafana_adaptive_traces_policy.example", "items.0.expires_at", "2099-12-31T23:59:59Z"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("grafana_adaptive_traces_policy.example", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "grafana_adaptive_traces_policy.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "grafana_adaptive_traces_policy" "example" {
  name       = "test_rule_updated"
  type       = "string_attribute"
  body       = jsonencode({
      key = "testfield.some-non-existing-field",
      values = [
          "ever"
      ]
  })
  expires_at = "2099-12-31T23:59:59Z"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("grafana_adaptive_traces_policy.example", "items.0.name", "test_rule_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
