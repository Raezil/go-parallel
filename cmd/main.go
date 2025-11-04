package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Raezil/go-parallel"
)

func main() {
	// IMPORTANT: Replace with your actual Parallel API key
	apiKey := "your-parallel-api-key"
	if apiKey == "your-parallel-api-key" {
		fmt.Println("Please replace 'your-parallel-api-key' with your actual Parallel API key.")
		return
	}

	// 1. Create a new client
	client := parallel.NewClient(apiKey)
	ctx := context.Background()

	// 2. Use the Search function
	fmt.Println("--- Testing Search ---")
	searchReq := parallel.ParallelSearchRequest{
		Objective:     "Find the latest news on AI.",
		SearchQueries: []string{"latest AI news"},
		MaxResults:    2,
	}
	searchResp, err := client.Search(ctx, searchReq)
	if err != nil {
		fmt.Printf("Error during Search: %v\n", err)
	} else {
		fmt.Printf("Search successful. Search ID: %s\n", searchResp.SearchID)
		for i, result := range searchResp.Results {
			fmt.Printf("  Result %d: %s (%s)\n", i+1, result.Title, result.URL)
		}
	}

	// 3. Use the Extract function
	fmt.Println("\n--- Testing Extract ---")
	extractReq := parallel.ParallelExtractRequest{
		URLs:      []string{"https://www.wired.com/category/artificial-intelligence/"},
		Objective: "Extract headlines from the page.",
	}
	extractResp, err := client.Extract(ctx, extractReq)
	if err != nil {
		fmt.Printf("Error during Extract: %v\n", err)
	} else {
		fmt.Printf("Extract successful. Extract ID: %s\n", extractResp.ExtractID)
		for i, result := range extractResp.Results {
			fmt.Printf("  Extracted Title %d: %s\n", i+1, result.Title)
		}
	}

	// 4. Use the RunTask, GetTask, and PollUntilComplete functions
	fmt.Println("\n--- Testing Tasks ---")
	taskReq := parallel.ParallelTaskRequest{
		Input:     "What were the key highlights of the latest Apple event?",
		Processor: "research-v1", // Example processor
	}
	taskResp, err := client.RunTask(ctx, taskReq)
	if err != nil {
		fmt.Printf("Error starting task: %v\n", err)
	} else {
		runID := taskResp.Output.RunID
		fmt.Printf("Task started successfully. Run ID: %s\n", runID)

		// 5. Get the initial status of the task
		fmt.Println("\n--- Testing GetTask ---")
		taskStatus, err := client.GetTask(ctx, runID)
		if err != nil {
			fmt.Printf("Error getting task status: %v\n", err)
		} else {
			fmt.Printf("Current task status: %s\n", taskStatus.Status)
		}

		// 6. Poll for the final result
		fmt.Println("\n--- Testing PollUntilComplete ---")
		fmt.Println("Polling for task completion (this might take a moment)...")
		finalResult, err := client.PollUntilComplete(ctx, runID, 5*time.Second)
		if err != nil {
			fmt.Printf("Error polling for task completion: %v\n", err)
		} else {
			fmt.Printf("Task finished with status: %s\n", finalResult.Status)
			if finalResult.Status == "completed" {
				fmt.Println("Task Output:", finalResult.Output)
			} else if finalResult.Error != nil {
				fmt.Println("Task Error:", finalResult.Error)
			}
		}
	}

	// 7. Use the Chat function
	fmt.Println("\n--- Testing Chat ---")
	chatReq := parallel.ParallelChatRequest{
		Model: "parallel-chat-v1",
		Messages: []parallel.ParallelChatMessage{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "What is the capital of France?"},
		},
	}
	chatResp, err := client.Chat(ctx, chatReq)
	if err != nil {
		fmt.Printf("Error during Chat: %v\n", err)
	} else {
		fmt.Println("Chat successful.")
		if len(chatResp.Choices) > 0 {
			fmt.Printf("  Assistant's response: %s\n", chatResp.Choices[0].Message.Content)
		}
	}
}
