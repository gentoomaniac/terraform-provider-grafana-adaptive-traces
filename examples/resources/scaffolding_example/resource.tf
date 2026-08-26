resource "grafana_adaptive_traces_policy" "example" {
  name       = "provider test rule"
  type       = "string_attribute"
  body       = jsonencode({
      key = "testfield.some-non-existing-field",
      values = [
          "ever"
      ]
  })
  expires_at = "2099-12-31T23:59:59Z"
}
