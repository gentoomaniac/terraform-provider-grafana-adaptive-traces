package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewAdaptiveTracesClient(tenantURL *url.URL, tenantId string, apiKey string) (*AdaptiveTracesClient, error) {
	if tenantURL == nil {
		return nil, fmt.Errorf("tenantUrl must be provided")
	}
	atc := AdaptiveTracesClient{
		client:    &http.Client{Timeout: 10 * time.Second},
		tenantURL: tenantURL,
		tenantID:  tenantId,
		apiKey:    apiKey,
	}

	return &atc, nil
}

type AdaptiveTracesClient struct {
	client    *http.Client
	tenantURL *url.URL
	tenantID  string
	apiKey    string
}

func (atc *AdaptiveTracesClient) sendRequest(ctx context.Context, verb string, endpoint string, payload []byte) (*http.Response, error) {
	path := atc.tenantURL.JoinPath("/adaptive-traces/api/v1", endpoint)

	var payloadReader io.Reader
	if len(payload) > 0 {
		payloadReader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, verb, path.String(), payloadReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(atc.tenantID, atc.apiKey)

	tflog.Debug(ctx, "sending request to grafana", map[string]any{"endpoint": endpoint, "full_path": path.String()})
	return atc.client.Do(req)
}

func (atc *AdaptiveTracesClient) GetPolicies(ctx context.Context) ([]Policy, error) {
	resp, err := atc.sendRequest(ctx, http.MethodGet, "policies", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("listing policies failed with status %d: %s", resp.StatusCode, string(body))
	}

	var policies []Policy
	if err := json.Unmarshal(body, &policies); err != nil {
		return nil, err
	}

	return policies, nil
}

func (atc *AdaptiveTracesClient) GetPolicy(ctx context.Context, id string) (*Policy, error) {
	resp, err := atc.sendRequest(ctx, http.MethodGet, path.Join("policies", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting policy `%s` failed with status %d: %s", id, resp.StatusCode, string(body))
	}

	var policies Policy
	if err := json.Unmarshal(body, &policies); err != nil {
		return nil, err
	}

	return &policies, nil
}

func (atc *AdaptiveTracesClient) CreatePolicy(ctx context.Context, policy Policy) (*Policy, error) {
	payload, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	resp, err := atc.sendRequest(ctx, http.MethodPost, "policies", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("creating policy failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createdPolicy Policy
	if err := json.Unmarshal(body, &createdPolicy); err != nil {
		return nil, err
	}

	return &createdPolicy, nil
}

func (atc *AdaptiveTracesClient) UpdatePolicy(ctx context.Context, policy Policy) (*Policy, error) {
	id := policy.ID

	// blank fields we don't want to send
	policy.ID = ""

	payload, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, string(payload))

	resp, err := atc.sendRequest(ctx, http.MethodPut, path.Join("policies", id), payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("updating policy failed with status %d: %s", resp.StatusCode, string(body))
	}

	var updatedPolicy Policy
	if err := json.Unmarshal(body, &updatedPolicy); err != nil {
		return nil, err
	}

	return &updatedPolicy, nil
}

func (atc *AdaptiveTracesClient) DeletePolicy(ctx context.Context, id string) error {
	resp, err := atc.sendRequest(ctx, http.MethodDelete, path.Join("policies", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("deleting policy failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createdPolicy Policy
	if err := json.Unmarshal(body, &createdPolicy); err != nil {
		return err
	}

	return nil
}
