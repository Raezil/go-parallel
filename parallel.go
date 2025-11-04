// path: parallel/parallel.go
package parallel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the core Parallel API client.
type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
	betaTag string
}

// NewClient creates a new Parallel API client with defaults.
func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://api.parallel.ai/v1beta",
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		betaTag: "search-extract-2025-10-10",
	}
}

// Search performs a semantic search query using the Parallel API.
func (c *Client) Search(ctx context.Context, req ParallelSearchRequest) (*ParallelSearchResponse, error) {
	url := fmt.Sprintf("%s/search", c.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("parallel-beta", c.betaTag)

	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API error: %s — %s", res.Status, string(body))
	}

	var out ParallelSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &out, nil
}

// RunTask launches a processing task (e.g., research, summarization, report generation).
func (c *Client) RunTask(ctx context.Context, req ParallelTaskRequest) (*ParallelTaskResponse, error) {
	url := fmt.Sprintf("%s/tasks/runs", c.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("parallel-beta", c.betaTag)

	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API error: %s — %s", res.Status, string(body))
	}

	var out ParallelTaskResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &out, nil
}

// GetTask retrieves the latest status or final output of a task run.
func (c *Client) GetTask(ctx context.Context, runID string) (*ParallelTaskResult, error) {
	url := fmt.Sprintf("%s/tasks/runs/%s", c.baseURL, runID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("parallel-beta", c.betaTag)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API error: %s — %s", res.Status, string(b))
	}

	var out ParallelTaskResult
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &out, nil
}

// PollUntilComplete continuously checks a task until its status is "completed" or context is canceled.
func (c *Client) PollUntilComplete(ctx context.Context, runID string, interval time.Duration) (*ParallelTaskResult, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			task, err := c.GetTask(ctx, runID)
			if err != nil {
				return nil, err
			}
			if task.Status == "completed" || task.Status == "failed" {
				return task, nil
			}
		}
	}
}

// Chat sends a chat completion request to Parallel's /chat/completions API.
func (c *Client) Chat(ctx context.Context, req ParallelChatRequest) (*ParallelChatResponse, error) {
	url := fmt.Sprintf("%s/chat/completions", c.baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API error: %s — %s", res.Status, string(b))
	}

	var out ParallelChatResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &out, nil
}

// Extract performs an extraction request on given URLs.
func (c *Client) Extract(ctx context.Context, req ParallelExtractRequest) (*ParallelExtractResponse, error) {
	url := fmt.Sprintf("%s/extract", c.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("parallel-beta", c.betaTag)

	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API error: %s — %s", res.Status, string(body))
	}

	var out ParallelExtractResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &out, nil
}
