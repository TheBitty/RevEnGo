package agent

// Options defines the configuration options for creating a new agent
type Options struct {
	// ModelName specifies which AI model to use (e.g., "deepseek:8b", "gemma3")
	ModelName string

	// OllamaEndpoint is the URL of the Ollama API endpoint
	OllamaEndpoint string

	// Verbose enables detailed logging when true
	Verbose bool

	// TrainingMode enables training capabilities when true
	TrainingMode bool

	// MaxTokens limits the maximum number of tokens in model responses
	MaxTokens int

	// Temperature controls randomness in model responses (0.0-1.0)
	Temperature float64

	// Concurrency controls the number of concurrent requests
	Concurrency int
}
