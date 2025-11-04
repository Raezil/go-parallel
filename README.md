# go-parallel

[![Go Reference](https://pkg.go.dev/badge/github.com/your-username/go-parallel.svg)](https://pkg.go.dev/github.com/your-username/go-parallel)

A Go client for the Parallel API, a powerful platform for web search, data extraction, and intelligent task automation.

## Overview

The Parallel API provides a suite of tools to build powerful applications:

*   **Search**: A semantic search engine to find the most relevant information on the web.
*   **Extract**: A tool to extract structured data from any URL.
*   **Tasks**: An engine to run complex, long-running tasks like market research, summarization, and report generation.
*   **Chat**: A conversational AI to power chatbots and other interactive experiences.

This Go client provides a simple and idiomatic way to interact with the Parallel API.

## Installation

```bash
go get github.com/your-username/go-parallel
```

*(Note: Replace `your-username` with the actual path when this is a real repository)*

## Usage

### Creating a Client

First, create a new client with your API key:

```go
import "github.com/your-username/go-parallel"

func main() {
    apiKey := "your-parallel-api-key"
    client := parallel.NewClient(apiKey)

    // ... use the client
}
```

### Search

Perform a semantic search:

```go
req := parallel.ParallelSearchRequest{
    Objective:     "Find information about the latest AI trends.",
    SearchQueries: []string{"AI trends 2025", "latest advancements in artificial intelligence"},
    MaxResults:    5,
}

resp, err := client.Search(context.Background(), req)
if err != nil {
    // handle error
}

fmt.Println("Search ID:", resp.SearchID)
for _, result := range resp.Results {
    fmt.Println("- ", result.Title)
}
```

### Extract

Extract content from URLs:

```go
req := parallel.ParallelExtractRequest{
    URLs:      []string{"https://www.example.com/article1", "https://www.another.com/news"},
    Objective: "Extract the main points from these articles.",
}

resp, err := client.Extract(context.Background(), req)
if err != nil {
    // handle error
}

fmt.Println("Extract ID:", resp.ExtractID)
for _, result := range resp.Results {
    fmt.Println("URL:", result.URL)
    fmt.Println("Title:", result.Title)
}
```

### Tasks

Run a task and poll for its completion:

```go
// Run the task
taskReq := parallel.ParallelTaskRequest{
    Input:     "Summarize the following text: ...",
    Processor: "summarize-v1",
}

taskResp, err := client.RunTask(context.Background(), taskReq)
if err != nil {
    // handle error
}

fmt.Println("Started task with Run ID:", taskResp.Output.RunID)

// Poll for completion
result, err := client.PollUntilComplete(context.Background(), taskResp.Output.RunID, 5*time.Second)
if err != nil {
    // handle error
}

if result.Status == "completed" {
    fmt.Println("Task completed successfully!")
    // Process result.Output
} else {
    fmt.Println("Task failed:", result.Error)
}
```

### Chat

Have a conversation with the Chat API:

```go
req := parallel.ParallelChatRequest{
    Model: "parallel-chat-v1",
    Messages: []parallel.ParallelChatMessage{
        {Role: "user", Content: "Hello, who are you?"},
    },
}

resp, err := client.Chat(context.Background(), req)
if err != nil {
    // handle error
}

fmt.Println("Response:", resp.Choices[0].Message.Content)
```

## API Documentation

For more detailed information about the API, see the [official Parallel API documentation](https://docs.parallel.ai/home).

## Testing

This project includes a suite of unit tests. To run the tests:

```bash
go test
```
