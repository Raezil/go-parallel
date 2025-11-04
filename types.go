// path: parallel/types.go
package parallel

import "time"

// ParallelSearchRequest defines the request body for Parallel API search.
type ParallelSearchRequest struct {
	Objective         string   `json:"objective"`
	SearchQueries     []string `json:"search_queries"`
	MaxResults        int      `json:"max_results"`
	MaxCharsPerResult int      `json:"max_chars_per_result"`
}

// ParallelSearchResponse is the full response from the API.
type ParallelSearchResponse struct {
	SearchID string           `json:"search_id"`
	Results  []ParallelResult `json:"results"`
}

// ParallelResult represents one search result.
type ParallelResult struct {
	URL      string   `json:"url"`
	Title    string   `json:"title"`
	Excerpts []string `json:"excerpts"`
}

// ParallelExtractRequest defines the request structure for /extract.
type ParallelExtractRequest struct {
	URLs        []string `json:"urls"`
	Objective   string   `json:"objective"`
	Excerpts    bool     `json:"excerpts"`
	FullContent bool     `json:"full_content"`
}

// ParallelExtractResponse represents the API’s extraction response.
type ParallelExtractResponse struct {
	ExtractID string             `json:"extract_id"`
	Results   []ParallelExtract  `json:"results"`
	Errors    []ParallelAPIError `json:"errors"`
}

// ParallelExtract represents a single extracted web page.
type ParallelExtract struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Excerpts    []string `json:"excerpts"`
	FullContent string   `json:"full_content"`
}

// ParallelAPIError captures any per-request errors.
type ParallelAPIError struct {
	Message string `json:"message"`
}

// ParallelTaskRequest defines the request structure for /tasks/runs.
type ParallelTaskRequest struct {
	Input     string `json:"input"`
	Processor string `json:"processor"`
}

// ParallelTaskResponse represents the structured output from Parallel’s task engine.
type ParallelTaskResponse struct {
	Output struct {
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
	} `json:"output"`
}

// ParallelTaskContent mirrors the “content” structure from your example.
type ParallelTaskContent struct {
	MarketSizeAndForecast struct {
		CAGR                string `json:"cagr"`
		MarketSegment       string `json:"market_segment"`
		CurrentValuation    string `json:"current_valuation"`
		ForecastedValuation string `json:"forecasted_valuation"`
		ForecastPeriod      string `json:"forecast_period"`
	} `json:"market_size_and_forecast"`

	CompanyProfiles []struct {
		CompanyName        string `json:"company_name"`
		StockTicker        string `json:"stock_ticker"`
		Revenue            string `json:"revenue"`
		MarketCap          string `json:"market_capitalization,omitempty"`
		MarketPosition     string `json:"market_position"`
		RecentDevelopments string `json:"recent_developments"`
	} `json:"company_profiles"`

	RecentMergersAndAcquisitions struct {
		AcquiringCompany string `json:"acquiring_company"`
		TargetCompany    string `json:"target_company"`
		DealSummary      string `json:"deal_summary"`
		Date             string `json:"date"`
	} `json:"recent_mergers_and_acquisitions"`

	GrowthOpportunities string `json:"growth_opportunities"`

	MarketSegmentationAnalysis struct {
		DominantSegment           string `json:"dominant_segment"`
		DominantSegmentShare      string `json:"dominant_segment_share"`
		FastestGrowingSegment     string `json:"fastest_growing_segment"`
		FastestGrowingSegmentCAGR string `json:"fastest_growing_segment_cagr"`
	} `json:"market_segmentation_analysis"`

	PubliclyTradedCompanies []struct {
		CompanyName string `json:"company_name"`
		StockTicker string `json:"stock_ticker"`
	} `json:"publicly_traded_hvac_companies"`
}

// ParallelBasis represents “basis” reasoning and citations.
type ParallelBasis struct {
	Field      string                  `json:"field"`
	Reasoning  string                  `json:"reasoning"`
	Citations  []ParallelBasisCitation `json:"citations"`
	Confidence string                  `json:"confidence"`
}

// ParallelBasisCitation provides structured citation data.
type ParallelBasisCitation struct {
	URL      string   `json:"url"`
	Excerpts []string `json:"excerpts"`
	Title    string   `json:"title"`
}

// ParallelTaskResult represents a **detailed** task status or completed output.
// This is what GetTask() and PollUntilComplete() return.
type ParallelTaskResult struct {
	RunID       string    `json:"run_id"`
	Status      string    `json:"status"`
	IsActive    bool      `json:"is_active"`
	Processor   string    `json:"processor"`
	Output      any       `json:"output"`       // raw JSON output from Parallel
	Error       any       `json:"error"`        // may contain error details or null
	Warnings    any       `json:"warnings"`     // optional warnings
	Metadata    any       `json:"metadata"`     // optional metadata
	TaskGroupID string    `json:"taskgroup_id"` // may be null
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
}

// ParallelChatRequest defines a chat completion request.
type ParallelChatRequest struct {
	Model          string                  `json:"model"`
	Messages       []ParallelChatMessage   `json:"messages"`
	Stream         bool                    `json:"stream"`
	ResponseFormat *ParallelResponseFormat `json:"response_format,omitempty"`
}

// ParallelChatMessage represents a chat message with a role and content.
type ParallelChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ParallelResponseFormat defines the structure for schema-based JSON output.
type ParallelResponseFormat struct {
	Type       string                         `json:"type"` // "json_schema"
	JSONSchema ParallelResponseJSONSchemaSpec `json:"json_schema"`
}

// ParallelResponseJSONSchemaSpec defines the schema metadata and body.
type ParallelResponseJSONSchemaSpec struct {
	Name   string         `json:"name"`
	Schema map[string]any `json:"schema"` // the full JSON schema
}

// ParallelChatResponse represents the API’s completion response.
type ParallelChatResponse struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Model   string               `json:"model"`
	Created int64                `json:"created"`
	Choices []ParallelChatChoice `json:"choices"`
	Usage   map[string]any       `json:"usage,omitempty"`
}

// ParallelChatChoice holds a single generated message.
type ParallelChatChoice struct {
	Index        int                 `json:"index"`
	Message      ParallelChatMessage `json:"message"`
	FinishReason string              `json:"finish_reason"`
}
