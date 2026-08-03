package hub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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

// PublishTemplate uploads template tar.gz archive and metadata to HarnessHub
func (c *Client) PublishTemplate(ctx context.Context, name, version, targetAgent, description string, archiveBytes []byte, blueprintJSON string) (*TemplateDetail, error) {
	if name == "" || version == "" {
		return nil, fmt.Errorf("name and version are required for publishing")
	}
	if targetAgent == "" {
		targetAgent = "antigravity"
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", name)
	_ = writer.WriteField("version", version)
	_ = writer.WriteField("target_agent", targetAgent)
	_ = writer.WriteField("description", description)
	_ = writer.WriteField("blueprint_json", blueprintJSON)

	part, err := writer.CreateFormFile("file", fmt.Sprintf("%s-%s.tar.gz", name, version))
	if err != nil {
		return nil, fmt.Errorf("failed to create multipart form file: %w", err)
	}
	if _, err := part.Write(archiveBytes); err != nil {
		return nil, fmt.Errorf("failed to write archive bytes into form: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	reqURL := fmt.Sprintf("%s/api/v1/templates", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to HarnessHub (%s): %w", c.BaseURL, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hub publish status %d: %s", resp.StatusCode, string(respBody))
	}

	var stdResp StandardResponse
	var detail TemplateDetail

	if err := json.Unmarshal(respBody, &stdResp); err == nil && stdResp.Data != nil {
		if !stdResp.Success && stdResp.Error != nil {
			return nil, fmt.Errorf("hub error [%s]: %s", stdResp.Error.Code, stdResp.Error.Message)
		}
		if err := json.Unmarshal(stdResp.Data, &detail); err == nil {
			return &detail, nil
		}
	}

	if err := json.Unmarshal(respBody, &detail); err != nil {
		return nil, fmt.Errorf("failed to decode response detail: %w", err)
	}

	return &detail, nil
}

// DeleteTemplate deletes a published template version from HarnessHub
func (c *Client) DeleteTemplate(ctx context.Context, name, version string) error {
	if name == "" || version == "" {
		return fmt.Errorf("name and version are required for deleting")
	}

	reqURL := fmt.Sprintf("%s/api/v1/templates/%s/%s", c.BaseURL, url.PathEscape(name), url.PathEscape(version))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to HarnessHub (%s): %w", c.BaseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("hub delete status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
