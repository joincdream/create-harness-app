package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// TemplateSummary represents item in templates list response
type TemplateSummary struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	TargetAgent string `json:"target_agent"`
	Description string `json:"description"`
}

// TemplateDetail represents detailed info with blueprint json
type TemplateDetail struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	TargetAgent   string `json:"target_agent"`
	Description   string `json:"description"`
	FilePath      string `json:"file_path"`
	FileSize      int64  `json:"file_size"`
	BlueprintJSON string `json:"blueprint_json"`
	PublisherID   int    `json:"publisher_id"`
}

// Client interacts with HarnessHub Backend REST API
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Client using env HARNESS_HUB_URL or default
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = os.Getenv("HARNESS_HUB_URL")
	}
	if baseURL == "" {
		baseURL = "http://localhost"
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// StandardResponse is the uniform API wrapper schema from Hub backend
type StandardResponse struct {
	Success   bool            `json:"success"`
	Data      json.RawMessage `json:"data"`
	Error     *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Timestamp string `json:"timestamp"`
}

// ListTemplates fetches templates from Hub backend
func (c *Client) ListTemplates(ctx context.Context, query, agent string) ([]TemplateSummary, error) {
	reqURL, err := url.Parse(c.BaseURL + "/api/v1/templates")
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	q := reqURL.Query()
	if query != "" {
		q.Set("q", query)
	}
	if agent != "" {
		q.Set("agent", agent)
	}
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to HarnessHub (%s): %w", c.BaseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hub returned status %d: %s", resp.StatusCode, string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var stdResp StandardResponse
	var list []TemplateSummary

	// Try decoding wrapped StandardResponse first
	if err := json.Unmarshal(bodyBytes, &stdResp); err == nil && stdResp.Data != nil {
		if !stdResp.Success && stdResp.Error != nil {
			return nil, fmt.Errorf("hub error [%s]: %s", stdResp.Error.Code, stdResp.Error.Message)
		}
		if err := json.Unmarshal(stdResp.Data, &list); err != nil {
			return nil, fmt.Errorf("failed to decode templates data array: %w", err)
		}
		return list, nil
	}

	// Fallback to raw array decoding
	if err := json.Unmarshal(bodyBytes, &list); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return list, nil
}

// GetTemplateDetail fetches single template detail from Hub backend
func (c *Client) GetTemplateDetail(ctx context.Context, name, version string) (*TemplateDetail, error) {
	if version == "" || version == "latest" {
		list, err := c.ListTemplates(ctx, name, "")
		if err == nil {
			for _, item := range list {
				if item.Name == name {
					version = item.Version
					break
				}
			}
		}
	}
	if version == "" {
		version = "v1.0.0"
	}
	reqURL := fmt.Sprintf("%s/api/v1/templates/%s/%s", c.BaseURL, url.PathEscape(name), url.PathEscape(version))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to HarnessHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hub returned status %d: %s", resp.StatusCode, string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var stdResp StandardResponse
	var detail TemplateDetail

	// Try decoding wrapped StandardResponse first
	if err := json.Unmarshal(bodyBytes, &stdResp); err == nil && stdResp.Data != nil {
		if !stdResp.Success && stdResp.Error != nil {
			return nil, fmt.Errorf("hub error [%s]: %s", stdResp.Error.Code, stdResp.Error.Message)
		}
		if err := json.Unmarshal(stdResp.Data, &detail); err != nil {
			return nil, fmt.Errorf("failed to decode template detail: %w", err)
		}
		return &detail, nil
	}

	// Fallback to raw detail decoding
	if err := json.Unmarshal(bodyBytes, &detail); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return &detail, nil
}

// DownloadTemplate downloads .tar.gz stream from Hub backend
func (c *Client) DownloadTemplate(ctx context.Context, name, version string) (io.ReadCloser, error) {
	if version == "" || version == "latest" {
		list, err := c.ListTemplates(ctx, name, "")
		if err == nil {
			for _, item := range list {
				if item.Name == name {
					version = item.Version
					break
				}
			}
		}
	}
	if version == "" {
		version = "v1.0.0"
	}
	reqURL := fmt.Sprintf("%s/api/v1/templates/%s/%s/download", c.BaseURL, url.PathEscape(name), url.PathEscape(version))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download from HarnessHub: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hub download status %d: %s", resp.StatusCode, string(body))
	}

	return resp.Body, nil
}
