package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	// Logger is the default logger for the parser package
	Logger = log.New(os.Stdout, "[PARSER] ", log.LstdFlags)
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

// Task represents a single task with its description
type Task struct {
	Task        string `json:"task"`
	Description string `json:"description"`
}

// OllamaResponse represents the response from the Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
}

// TaskParser handles the parsing of user input into structured tasks using Ollama
type TaskParser struct {
	modelName string
	client    *http.Client
}

// NewTaskParser creates a new TaskParser instance
func NewTaskParser(modelName string) *TaskParser {
	LogInfo("Creating new task parser with model: %s", modelName)
	return &TaskParser{
		modelName: modelName,
		client:    &http.Client{},
	}
}

// GetModelName returns the name of the model being used
func (p *TaskParser) GetModelName() string {
	return p.modelName
}

// SendRequest sends a request to Ollama and returns the response
func (p *TaskParser) SendRequest(requestBody map[string]interface{}) (string, error) {
	LogDebug("Sending request to Ollama: %v", requestBody)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		LogError("Error marshaling request: %v", err)
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := p.client.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		LogError("Error sending request to Ollama: %v", err)
		return "", fmt.Errorf("error sending request to Ollama: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogError("Error reading response body: %v", err)
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		LogError("Error parsing response: %v", err)
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	responseText, ok := response["response"].(string)
	if !ok {
		LogError("Response does not contain 'response' field")
		return "", fmt.Errorf("response does not contain 'response' field")
	}

	return responseText, nil
}

// ParseInput takes user input and returns a list of structured tasks
func (p *TaskParser) ParseInput(userInput string) ([]Task, error) {
	LogInfo("Parsing input: %s", userInput)

	prompt := fmt.Sprintf(`Given the following user input, break it down into a list of specific tasks or instructions.
Format each task with a 'task' field and a 'description' field.
Return the response as a valid JSON array.

User Input: %s

Example format:
[
    {"task": "Task 1", "description": "Description of task 1"},
    {"task": "Task 2", "description": "Description of task 2"}
]`, userInput)

	requestBody := map[string]interface{}{
		"model":  p.modelName,
		"prompt": prompt,
		"stream": false,
	}

	response, err := p.SendRequest(requestBody)
	if err != nil {
		LogError("Error getting response from Ollama: %v", err)
		return nil, fmt.Errorf("error getting response from Ollama: %v", err)
	}

	// Clean up the response to ensure it's valid JSON
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	// Extract JSON array from the response
	startIdx := strings.Index(response, "[")
	endIdx := strings.LastIndex(response, "]") + 1
	if startIdx == -1 || endIdx == 0 {
		return nil, fmt.Errorf("could not find JSON array in response")
	}

	jsonStr := response[startIdx:endIdx]

	// Clean up the JSON string by removing any problematic characters
	jsonStr = strings.ReplaceAll(jsonStr, "\n", " ")
	jsonStr = strings.ReplaceAll(jsonStr, "\r", " ")
	jsonStr = strings.ReplaceAll(jsonStr, "\t", " ")

	var tasks []Task
	if err := json.Unmarshal([]byte(jsonStr), &tasks); err != nil {
		return nil, fmt.Errorf("error parsing tasks: %v", err)
	}

	LogInfo("Parsed input into %d tasks", len(tasks))
	return tasks, nil
}
