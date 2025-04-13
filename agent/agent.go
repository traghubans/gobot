package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// Default Ollama endpoint if not specified in environment
	DefaultOllamaHost = "http://localhost:11434"

	// Status constants
	StatusStarted   = "started"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

var (
	// Logger is the default logger for the agent package
	Logger = log.New(os.Stdout, "[AGENT] ", log.LstdFlags)
)

// LogInfo logs an informational message
func LogInfo(format string, args ...interface{}) {
	Logger.Printf("[INFO] "+format, args...)
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	Logger.Printf("[ERROR] "+format, args...)
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	Logger.Printf("[DEBUG] "+format, args...)
}

// QueryResult represents the result of a query
type QueryResult struct {
	Query     string
	Answer    string
	Timestamp time.Time
	Status    string
	Error     error
}

// Agent represents an AI agent that can execute queries and maintain context
type Agent struct {
	modelName string
	client    *http.Client
	context   []QueryResult
}

// NewAgent creates a new Agent instance
func NewAgent(modelName string) *Agent {
	if modelName == "" {
		modelName = "mistral"
	}

	LogInfo("Creating new agent with model: %s", modelName)
	return &Agent{
		modelName: modelName,
		client:    &http.Client{},
		context:   make([]QueryResult, 0),
	}
}

// ProcessQuery processes a user query and returns a response with context
func (a *Agent) ProcessQuery(query string) (QueryResult, error) {
	if query == "" {
		return QueryResult{}, fmt.Errorf("query cannot be empty")
	}

	LogInfo("Processing query: %s", query)
	result := QueryResult{
		Query:     query,
		Timestamp: time.Now(),
		Status:    StatusStarted,
	}

	response, err := a.sendRequest(a.buildPromptWithContext(query))
	if err != nil {
		LogError("Error getting response: %v", err)
		result.Status = StatusFailed
		result.Error = err
		result.Answer = fmt.Sprintf("Error: %v", err)
		a.context = append(a.context, result)
		return result, err
	}

	result.Status = StatusCompleted
	result.Answer = response
	a.context = append(a.context, result)
	LogInfo("Query processed successfully")

	return result, nil
}

// buildPromptWithContext creates a prompt that includes context from previous queries
func (a *Agent) buildPromptWithContext(query string) string {
	var prompt strings.Builder

	prompt.WriteString("You are a helpful AI assistant. Please format your responses according to these rules:\n")
	prompt.WriteString("1. If providing a list of items, use numbered points (1., 2., etc.) with each point on a new line\n")
	prompt.WriteString("2. If explaining concepts, separate different points with line breaks\n")
	prompt.WriteString("3. For any lists or steps, add a line break before and after the list\n")
	prompt.WriteString("4. Keep paragraphs concise and separated by line breaks\n")
	prompt.WriteString("5. Use bullet points (•) for sub-items or related points\n\n")

	if len(a.context) > 0 {
		prompt.WriteString("Previous conversation context:\n")
		for i, prev := range a.context {
			if prev.Status == StatusCompleted {
				prompt.WriteString(fmt.Sprintf("%d. Q: %s\n", i+1, prev.Query))
				prompt.WriteString(fmt.Sprintf("   A: %s\n\n", prev.Answer))
			}
		}
	}

	prompt.WriteString("Current query: ")
	prompt.WriteString(query)
	prompt.WriteString("\n\nPlease provide a helpful response that builds on the previous context if relevant. Remember to format your response according to the rules above.")

	return prompt.String()
}

// getOllamaEndpoint returns the Ollama API endpoint using environment variable if set
func getOllamaEndpoint() string {
	if host := os.Getenv("OLLAMA_HOST"); host != "" {
		return fmt.Sprintf("http://%s:11434/api/generate", host)
	}
	return DefaultOllamaHost + "/api/generate"
}

// sendRequest sends a request to Ollama and returns the response
func (a *Agent) sendRequest(prompt string) (string, error) {
	LogDebug("Sending request to Ollama")

	requestBody := map[string]interface{}{
		"model":  a.modelName,
		"prompt": prompt,
		"stream": false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := a.client.Post(getOllamaEndpoint(), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error sending request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	responseText, ok := response["response"].(string)
	if !ok {
		return "", fmt.Errorf("response does not contain 'response' field")
	}

	return responseText, nil
}

// GetContext returns the current context of queries
func (a *Agent) GetContext() []QueryResult {
	return a.context
}

// ClearContext clears the agent's context
func (a *Agent) ClearContext() {
	a.context = make([]QueryResult, 0)
}

// Close closes the agent and cleans up resources
func (a *Agent) Close() error {
	a.ClearContext()
	return nil
}
