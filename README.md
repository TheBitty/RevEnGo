# RevEnGo

RevEnGo is a specialized note-taking application designed specifically for reverse engineering tasks. It provides a clean, intuitive interface for organizing and documenting findings during the reverse engineering process, helping analysts keep track of their discoveries and insights.

## Key Features

- **Purpose-Built Interface**: Designed specifically for reverse engineering workflows
- **Structured Organization**: Organize notes by projects and tags for easy retrieval
- **Persistent Storage**: All notes are automatically saved for later reference
- **Dark Theme**: Easy on the eyes during long analysis sessions
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Technical Overview

RevEnGo is built using Go and the [Fyne](https://fyne.io/) UI toolkit, providing a lightweight, native-feeling application across all supported platforms. The application uses a simple, file-based storage system that saves notes as JSON files, making them easy to back up or version control.

### Architecture

RevEnGo follows a clean separation of concerns with the following components:

- **Models**: Data structures and storage interfaces for notes and projects
- **UI Components**: Reusable UI elements that make up the application interface
- **Main Application**: Initialization, setup, and assembly of the application

## Requirements

- Go 1.18 or later
- [Fyne toolkit dependencies](https://developer.fyne.io/started/)

### Fyne Dependencies

Fyne requires certain system dependencies to build and run GUI applications:

- **macOS**: Xcode (or Command Line Tools) - `xcode-select --install`
- **Windows**: GCC (via MSYS2 or MinGW) and a C compiler
- **Linux**: GCC, X11 and GL development libraries

## Installation

```bash
# Clone the repository
git clone https://github.com/leog/RevEnGo.git
cd RevEnGo

# Install dependencies and build
go mod tidy
go build -o revengo

# Run the application
./revengo
```

## Usage Guide

### Creating and Managing Notes

1. **Create a New Note**: Click the "New Note" button in the header
2. **Set a Title**: Enter a descriptive title for your note
3. **Write Content**: Document your reverse engineering findings in the main content area
4. **Add Tags**: Use tags to categorize your notes (e.g., "buffer-overflow", "x86", "encryption")
5. **Save**: Click the "Save" button to store your note

### Organizing Notes

- **Projects**: Group related notes under projects for better organization
- **Tags**: Use tags to create cross-cutting categories across projects
- **Search**: Find notes quickly using the search bar

## Project Structure

```
RevEnGo/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── internal/               # Internal application code
│   ├── models/             # Data models
│   │   ├── note.go         # Note data model and storage
│   │   └── project.go      # Project data model and storage
│   └── ui/                 # User interface components
│       └── components/     # Reusable UI elements
│           ├── header.go   # Application header
│           ├── notepad.go  # Note editing component
│           └── sidebar.go  # Navigation sidebar
└── pkg/                    # Public libraries (future expansion)
```

## Future Enhancements

- **Binary Analysis Integration**: Direct integration with binary analysis tools
- **Metadata Extraction**: Automatic extraction of functions, strings, and other binary metadata
- **Memory Visualization**: Tools for visualizing memory structures and layouts
- **Disassembler Integration**: Connect with popular disassemblers for seamless workflow
- **Collaboration Features**: Share and collaborate on reverse engineering notes

## Contributing

Contributions to RevEnGo are welcome! Please feel free to submit issues or pull requests.

1. Fork the repository
2. Create a feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request

## License

[MIT License](LICENSE)

## Acknowledgments

- The [Fyne](https://fyne.io/) team for their excellent GUI toolkit
- The Go community for the powerful and elegant Go programming language
