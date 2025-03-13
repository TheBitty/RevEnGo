package agent

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/yourusername/RevEnGo/internal/models"
)

// DefaultOllamaEndpoint is the default Ollama API endpoint
const DefaultOllamaEndpoint = "http://localhost:11434/api"

// Agent represents the AI-powered reverse engineering agent
type Agent struct {
	options Options
	model   models.Model
	mutex   sync.RWMutex
	logger  *Logger
}

// Logger provides a simple logging interface
type Logger struct {
	verbose bool
}

// Log logs a message if verbose mode is enabled
func (l *Logger) Log(format string, args ...interface{}) {
	if l.verbose {
		fmt.Printf("[RevEnGo] "+format+"\n", args...)
	}
}

// NewAgent creates a new Agent with the given options
func NewAgent(options Options) (*Agent, error) {
	// Apply default values for unspecified options
	if options.ModelName == "" {
		options.ModelName = "deepseek:8b"
	}

	if options.OllamaEndpoint == "" {
		options.OllamaEndpoint = DefaultOllamaEndpoint
	}

	if options.MaxTokens == 0 {
		options.MaxTokens = 1024
	}

	if options.Temperature == 0 {
		options.Temperature = 0.7
	}

	if options.Concurrency == 0 {
		options.Concurrency = 2
	}

	// Create the logger
	logger := &Logger{
		verbose: options.Verbose,
	}

	// Initialize the agent
	agent := &Agent{
		options: options,
		logger:  logger,
	}

	// Initialize the appropriate model
	var err error
	switch options.ModelName {
	case "deepseek:8b":
		agent.model, err = models.NewDeepSeekModel(options.OllamaEndpoint, options.MaxTokens, options.Temperature)
	case "gemma3":
		agent.model, err = models.NewGemmaModel(options.OllamaEndpoint, options.MaxTokens, options.Temperature)
	default:
		return nil, fmt.Errorf("unsupported model: %s", options.ModelName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize model: %w", err)
	}

	agent.logger.Log("Agent initialized with model %s", options.ModelName)
	return agent, nil
}

// AnalysisResult contains the results of a file analysis
type AnalysisResult struct {
	Filename        string
	FileType        string
	FileSize        int64
	Findings        []Finding
	Summary         string
	Vulnerabilities []Vulnerability
}

// Finding represents a notable item found during analysis
type Finding struct {
	Type        string
	Description string
	Location    string
	Severity    string
}

// Vulnerability represents a potential security vulnerability
type Vulnerability struct {
	Type        string
	Description string
	Location    string
	Severity    string
	CVSS        float64
	Remediation string
}

// AnalyzeFile performs AI-powered analysis on the given file
func (a *Agent) AnalyzeFile(filePath string) (*AnalysisResult, error) {
	a.logger.Log("Analyzing file: %s", filePath)

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to access file: %w", err)
	}

	// Read file content (with size limits for safety)
	const maxSize = 10 * 1024 * 1024 // 10MB limit
	if fileInfo.Size() > maxSize {
		return nil, errors.New("file too large for analysis (max 10MB)")
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Create analysis prompts based on file type
	fileExt := filepath.Ext(filePath)
	fileType := determineFileType(fileExt, fileContent)

	// Initialize result
	result := &AnalysisResult{
		Filename: filepath.Base(filePath),
		FileType: fileType,
		FileSize: fileInfo.Size(),
		Findings: []Finding{},
	}

	// Process the file with concurrent analysis tasks
	a.logger.Log("File type detected: %s", fileType)

	// Create a wait group for concurrent processing
	var wg sync.WaitGroup

	// Use a channel to collect findings
	findingsChan := make(chan Finding, 10)
	vulnChan := make(chan Vulnerability, 10)
	summaryChan := make(chan string, 1)

	// Launch concurrent analyses
	wg.Add(3)

	// Task 1: Basic information extraction
	go func() {
		defer wg.Done()
		prompt := generateInfoExtractionPrompt(fileType, fileContent)
		response, err := a.model.Generate(prompt)
		if err != nil {
			a.logger.Log("Error in info extraction: %v", err)
			return
		}

		findings := parseFindings(response)
		for _, f := range findings {
			findingsChan <- f
		}
	}()

	// Task 2: Vulnerability analysis
	go func() {
		defer wg.Done()
		prompt := generateVulnerabilityPrompt(fileType, fileContent)
		response, err := a.model.Generate(prompt)
		if err != nil {
			a.logger.Log("Error in vulnerability analysis: %v", err)
			return
		}

		vulns := parseVulnerabilities(response)
		for _, v := range vulns {
			vulnChan <- v
		}
	}()

	// Task 3: Generate summary
	go func() {
		defer wg.Done()
		prompt := generateSummaryPrompt(fileType, fileContent)
		response, err := a.model.Generate(prompt)
		if err != nil {
			a.logger.Log("Error in summary generation: %v", err)
			return
		}

		summaryChan <- response
	}()

	// Wait for all tasks to complete
	go func() {
		wg.Wait()
		close(findingsChan)
		close(vulnChan)
		close(summaryChan)
	}()

	// Collect results
	for finding := range findingsChan {
		result.Findings = append(result.Findings, finding)
	}

	for vuln := range vulnChan {
		result.Vulnerabilities = append(result.Vulnerabilities, vuln)
	}

	// Get summary
	for summary := range summaryChan {
		result.Summary = summary
	}

	a.logger.Log("Analysis completed with %d findings and %d vulnerabilities",
		len(result.Findings), len(result.Vulnerabilities))

	return result, nil
}

// Train trains the model using the provided dataset
func (a *Agent) Train(datasetPath, outputPath string) error {
	if !a.options.TrainingMode {
		return errors.New("agent not initialized in training mode")
	}

	a.logger.Log("Starting training with dataset: %s", datasetPath)

	// Check if dataset exists
	_, err := os.Stat(datasetPath)
	if err != nil {
		return fmt.Errorf("failed to access dataset: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Delegate training to the model
	if trainer, ok := a.model.(models.Trainable); ok {
		err = trainer.Train(datasetPath, outputPath)
		if err != nil {
			return fmt.Errorf("training failed: %w", err)
		}
	} else {
		return fmt.Errorf("model %s does not support training", a.options.ModelName)
	}

	a.logger.Log("Training completed successfully")
	return nil
}

// Helper functions

// determineFileType attempts to identify the type of file based on extension and content
func determineFileType(ext string, content []byte) string {
	// Simple file type detection based on extension
	switch ext {
	case ".exe", ".dll":
		return "Windows PE Executable"
	case ".elf", ".so":
		return "ELF Binary"
	case ".jar":
		return "Java Archive"
	case ".class":
		return "Java Bytecode"
	case ".js":
		return "JavaScript"
	case ".py":
		return "Python"
	case ".go":
		return "Go"
	case ".c", ".cpp", ".h", ".hpp":
		return "C/C++"
	default:
		// Try to detect binary vs text
		if isBinary(content) {
			return "Binary"
		}
		return "Text"
	}
}

// isBinary does a simple check to determine if content is likely binary
func isBinary(content []byte) bool {
	// Check for NULL bytes which are common in binary files
	for _, b := range content[:min(len(content), 1000)] {
		if b == 0 {
			return true
		}
	}
	return false
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// generateInfoExtractionPrompt creates a prompt for extracting basic information
func generateInfoExtractionPrompt(fileType string, content []byte) string {
	return fmt.Sprintf(`Analyze this %s file and extract key information:
Content sample: %s
Provide detailed findings about the structure, imports, dependencies, or other notable elements.
Format your response in a structured way, one finding per line, with the format:
TYPE: DESCRIPTION: LOCATION: SEVERITY
`,
		fileType,
		string(content[:min(len(content), 1000)]))
}

// generateVulnerabilityPrompt creates a prompt for vulnerability analysis
func generateVulnerabilityPrompt(fileType string, content []byte) string {
	return fmt.Sprintf(`Analyze this %s file for security vulnerabilities:
Content sample: %s
Look for common issues like memory safety, input validation, insecure functions, etc.
Format your response in a structured way, one vulnerability per line, with the format:
TYPE: DESCRIPTION: LOCATION: SEVERITY: CVSS: REMEDIATION
`,
		fileType,
		string(content[:min(len(content), 1000)]))
}

// generateSummaryPrompt creates a prompt for generating a summary
func generateSummaryPrompt(fileType string, content []byte) string {
	return fmt.Sprintf(`Provide a concise summary of this %s file:
Content sample: %s
What is its likely purpose? What are its main components? Is it potentially malicious?
Provide your analysis in a paragraph form.
`,
		fileType,
		string(content[:min(len(content), 1000)]))
}

// parseFindings parses the model response into Finding structures
func parseFindings(response string) []Finding {
	// Simplified parsing for example purposes
	// In a real application, use proper parsing logic based on the expected format
	findings := []Finding{}

	// Add a sample finding for demonstration
	findings = append(findings, Finding{
		Type:        "Sample",
		Description: "This is a sample finding from the response: " + response[:min(len(response), 100)],
		Location:    "N/A",
		Severity:    "Low",
	})

	return findings
}

// parseVulnerabilities parses the model response into Vulnerability structures
func parseVulnerabilities(response string) []Vulnerability {
	// Simplified parsing for example purposes
	// In a real application, use proper parsing logic based on the expected format
	vulns := []Vulnerability{}

	// Add a sample vulnerability for demonstration
	vulns = append(vulns, Vulnerability{
		Type:        "Sample",
		Description: "This is a sample vulnerability from the response: " + response[:min(len(response), 100)],
		Location:    "N/A",
		Severity:    "Low",
		CVSS:        3.2,
		Remediation: "Example remediation steps would go here.",
	})

	return vulns
}
