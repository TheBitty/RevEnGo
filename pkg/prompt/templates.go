package prompt

import (
	"bytes"
	"fmt"
	"text/template"
)

// Template types for different reverse engineering tasks
const (
	BinaryAnalysisTemplate     = "binary_analysis"
	VulnerabilityTemplate      = "vulnerability_analysis"
	SourceCodeAnalysisTemplate = "source_code_analysis"
	DecompilationTemplate      = "decompilation"
	StringAnalysisTemplate     = "string_analysis"
)

// TemplateData contains data for filling a prompt template
type TemplateData struct {
	// Any content that needs to be analyzed
	Content string

	// Type of file being analyzed
	FileType string

	// Architecture of the binary (x86, x86_64, ARM, etc.)
	Architecture string

	// Format preferences for the response
	OutputFormat string

	// Level of detail requested (brief, detailed, etc.)
	DetailLevel string

	// Focus area for the analysis (e.g., "security", "performance")
	Focus string

	// Any custom values to be used in templates
	Custom map[string]interface{}
}

// GetTemplateByName returns the template text for a given template name
func GetTemplateByName(name string) string {
	switch name {
	case BinaryAnalysisTemplate:
		return binaryAnalysisTemplate
	case VulnerabilityTemplate:
		return vulnerabilityTemplate
	case SourceCodeAnalysisTemplate:
		return sourceCodeAnalysisTemplate
	case DecompilationTemplate:
		return decompilationTemplate
	case StringAnalysisTemplate:
		return stringAnalysisTemplate
	default:
		return genericTemplate
	}
}

// Format fills a template with the provided data
func Format(templateName string, data TemplateData) (string, error) {
	// Get the template text
	templateText := GetTemplateByName(templateName)

	// Parse the template
	tmpl, err := template.New(templateName).Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Fill the template
	var result bytes.Buffer
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// Template definitions

// genericTemplate is a generic template for any task
const genericTemplate = `
Analyze the following content:

{{.Content}}

Provide a {{.DetailLevel}} analysis focusing on {{.Focus}}.
Output your response in {{.OutputFormat}} format.
`

// binaryAnalysisTemplate is a template for binary analysis
const binaryAnalysisTemplate = `
Analyze the following binary file:

File Type: {{.FileType}}
Architecture: {{.Architecture}}

Content:
{{.Content}}

Provide a {{.DetailLevel}} analysis focusing on:
1. Binary structure and sections
2. Identified functions and their purpose
3. Potential vulnerabilities or security issues
4. Notable strings and data found in the binary
5. Recommended further analysis steps

Output your response in {{.OutputFormat}} format.
`

// vulnerabilityTemplate is a template for vulnerability analysis
const vulnerabilityTemplate = `
Analyze the following code or binary for security vulnerabilities:

File Type: {{.FileType}}
{{if .Architecture}}Architecture: {{.Architecture}}{{end}}

Content:
{{.Content}}

Provide a {{.DetailLevel}} analysis of potential security vulnerabilities, including:
1. Memory safety issues (buffer overflows, use-after-free, etc.)
2. Input validation issues
3. Authentication/authorization flaws
4. Insecure cryptographic implementations
5. Race conditions
6. Injection vulnerabilities
7. Severity assessment (CVSS if possible)
8. Remediation recommendations

Output your response in {{.OutputFormat}} format.
`

// sourceCodeAnalysisTemplate is a template for source code analysis
const sourceCodeAnalysisTemplate = `
Analyze the following source code:

File Type: {{.FileType}}

Content:
{{.Content}}

Provide a {{.DetailLevel}} analysis focusing on:
1. Program functionality and purpose
2. Code structure and design patterns
3. Potential bugs or logic errors
4. Security considerations
5. Performance implications
6. Code quality assessment
7. Recommendations for improvement

Output your response in {{.OutputFormat}} format.
`

// decompilationTemplate is a template for decompiling binary code
const decompilationTemplate = `
Decompile the following binary code:

Architecture: {{.Architecture}}

Content:
{{.Content}}

Provide the closest possible source code representation of this binary, focusing on:
1. Function boundaries and calling conventions
2. Control flow structures (if/else, loops)
3. Variable usage and data structures
4. Comments explaining unclear sections
5. Potential purpose of the code

Output your response as clean, readable source code with appropriate comments.
`

// stringAnalysisTemplate is a template for analyzing strings
const stringAnalysisTemplate = `
Analyze the following strings extracted from a binary:

File Type: {{.FileType}}
{{if .Architecture}}Architecture: {{.Architecture}}{{end}}

Strings:
{{.Content}}

Provide a {{.DetailLevel}} analysis focusing on:
1. Interesting strings that reveal functionality
2. File paths, URLs, or network indicators
3. Error messages and what they indicate
4. Potential command strings or arguments
5. Cryptographic artifacts (keys, certificates)
6. Evidence of malicious behavior
7. Likely purpose of the binary based on strings

Output your response in {{.OutputFormat}} format.
`
