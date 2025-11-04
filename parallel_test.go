
package parallel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey to be %s, got %s", apiKey, client.apiKey)
	}

	if client.baseURL != "https://api.parallel.ai/v1beta" {
		t.Errorf("Expected baseURL to be https://api.parallel.ai/v1beta, got %s", client.baseURL)
	}

	if client.betaTag != "search-extract-2025-10-10" {
		t.Errorf("Expected betaTag to be search-extract-2025-10-10, got %s", client.betaTag)
	}
}

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1beta/search" {
			t.Errorf("Expected to request '/v1beta/search', got %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "test-api-key" {
			t.Errorf("Expected x-api-key header to be 'test-api-key', got %s", r.Header.Get("x-api-key"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ParallelSearchResponse{
			SearchID: "test-search-id",
			Results: []ParallelResult{
				{
					URL:   "https://example.com",
					Title: "Test Title",
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL + "/v1beta"

	req := ParallelSearchRequest{
		Objective: "test objective",
	}

	resp, err := client.Search(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.SearchID != "test-search-id" {
		t.Errorf("Expected SearchID to be 'test-search-id', got %s", resp.SearchID)
	}

	if len(resp.Results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(resp.Results))
	}

	if resp.Results[0].Title != "Test Title" {
		t.Errorf("Expected result title to be 'Test Title', got %s", resp.Results[0].Title)
	}
}

func TestExtract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1beta/extract" {
			t.Errorf("Expected to request '/v1beta/extract', got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ParallelExtractResponse{
			ExtractID: "test-extract-id",
			Results: []ParallelExtract{
				{
					URL:   "https://example.com",
					Title: "Test Extract Title",
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL + "/v1beta"

	req := ParallelExtractRequest{
		URLs: []string{"https://example.com"},
	}

	resp, err := client.Extract(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.ExtractID != "test-extract-id" {
		t.Errorf("Expected ExtractID to be 'test-extract-id', got %s", resp.ExtractID)
	}

	if len(resp.Results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(resp.Results))
	}

	if resp.Results[0].Title != "Test Extract Title" {
		t.Errorf("Expected result title to be 'Test Extract Title', got %s", resp.Results[0].Title)
	}
}

func TestRunTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1beta/tasks/runs" {
			t.Errorf("Expected to request '/v1beta/tasks/runs', got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ParallelTaskResponse{
			Output: struct {
				Content     ParallelTaskContent `json:"content"`
				Basis       []ParallelBasis     `json:"basis"`
				RunID       string              `json:"run_id"`
				Status      string              `json:"status"`
				CreatedAt   time.Time           `json:"created_at"`
				CompletedAt time.Time           `json:"completed_at"`
				Processor   string              `json:"processor"`
				Warnings    any                 `json:"warnings"`
				Error       any                 `json:"error"`
				TaskGroup   any                 `json:"taskgroup_id"`
			}{
				RunID: "test-run-id",
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL + "/v1beta"

	req := ParallelTaskRequest{
		Input: "test input",
	}

	resp, err := client.RunTask(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Output.RunID != "test-run-id" {
		t.Errorf("Expected RunID to be 'test-run-id', got %s", resp.Output.RunID)
	}
}

func TestGetTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1beta/tasks/runs/test-run-id" {
			t.Errorf("Expected to request '/v1beta/tasks/runs/test-run-id', got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ParallelTaskResult{
			RunID:  "test-run-id",
			Status: "completed",
		})
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL + "/v1beta"

	resp, err := client.GetTask(context.Background(), "test-run-id")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.RunID != "test-run-id" {
		t.Errorf("Expected RunID to be 'test-run-id', got %s", resp.RunID)
	}
	if resp.Status != "completed" {
		t.Errorf("Expected Status to be 'completed', got %s", resp.Status)
	}
}

func TestPollUntilComplete(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		status := "running"
		if callCount > 2 {
			status = "completed"
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ParallelTaskResult{
			RunID:  "test-run-id",
			Status: status,
		})
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL + "/v1beta"

	resp, err := client.PollUntilComplete(context.Background(), "test-run-id", 10*time.Millisecond)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Status != "completed" {
		t.Errorf("Expected final status to be 'completed', got %s", resp.Status)
	}
	if callCount != 3 {
		t.Errorf("Expected GetTask to be called 3 times, got %d", callCount)
	}
}

func TestChat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1beta/chat/completions" {
			t.Errorf("Expected to request '/v1beta/chat/completions', got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header to be 'Bearer test-api-key', got %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ParallelChatResponse{
			ID: "test-chat-id",
			Choices: []ParallelChatChoice{
				{
					Message: ParallelChatMessage{
						Role:    "assistant",
						Content: "Hello!",
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL + "/v1beta"

	req := ParallelChatRequest{
		Model: "test-model",
		Messages: []ParallelChatMessage{
			{
				Role:    "user",
				Content: "Hi",
			},
		},
	}

	resp, err := client.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.ID != "test-chat-id" {
		t.Errorf("Expected ID to be 'test-chat-id', got %s", resp.ID)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("Expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Content != "Hello!" {
		t.Errorf("Expected message content to be 'Hello!', got %s", resp.Choices[0].Message.Content)
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "invalid api key")
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL + "/v1beta"

	_, err := client.Search(context.Background(), ParallelSearchRequest{})
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	expectedError := "API error: 401 Unauthorized — invalid api key"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}
