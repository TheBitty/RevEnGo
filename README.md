# RevEnGo

A high-performance Go-based reverse engineering automation tool designed to analyze any type of files. RevEngGo provides automated analysis for finding strings, detecting buffer overflow vulnerabilities, it mainly performs static analysis.

## Overview

RevEnGo is an AI-powered agent for automating reverse engineering tasks. It leverages large language models (DeepSeek:8b and Gemma3) through OLLAMA to provide intelligent analysis of binaries, source code, and other artifacts.

## Features

- Automated binary and source code analysis
- AI-powered reverse engineering using DeepSeek:8b and Gemma3 models
- Custom dataset training to enhance AI capabilities
- Concurrent processing for high-performance analysis
- Modular architecture for extensibility

## Requirements

- Go 1.24 or later
- OLLAMA installed and configured
- Access to DeepSeek:8b and Gemma3 models

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/RevEnGo.git
cd RevEnGo

# Build the project
go build -o revengo ./cmd/revengo
```

## Usage

```bash
# Run basic reverse engineering on a binary
./revengo analyze -file /path/to/binary

# Train on custom dataset
./revengo train -dataset /path/to/dataset -output /path/to/output
```

## Architecture

RevEnGo is structured with a clear separation of concerns:

- `cmd/` - Application entry points
- `internal/` - Private application code
  - `agent/` - Core agent implementation
  - `models/` - AI model integrations
  - `training/` - Dataset training utilities
  - `reverse/` - Reverse engineering tools
- `pkg/` - Public libraries that can be used by external applications

## License

[MIT License](LICENSE)
