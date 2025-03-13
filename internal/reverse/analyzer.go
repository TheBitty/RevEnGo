package reverse

import (
	"bytes"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// FileInfo contains basic information about an analyzed file
type FileInfo struct {
	Path       string
	Name       string
	Size       int64
	Type       string
	Arch       string
	Imports    []string
	Exports    []string
	Strings    []string
	Sections   []Section
	IsStripped bool
}

// Section represents a section in a binary file
type Section struct {
	Name         string
	Addr         uint64
	Size         uint64
	Offset       uint64
	Flags        string
	IsExecutable bool
}

// AnalyzeFile performs static analysis on a file and returns information about it
func AnalyzeFile(path string) (*FileInfo, error) {
	// Check if file exists
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access file: %w", err)
	}

	// Create result structure
	result := &FileInfo{
		Path: path,
		Name: filepath.Base(path),
		Size: fileInfo.Size(),
	}

	// Try to identify file type
	fileType, err := identifyFileType(path)
	if err != nil {
		// Not a recognized binary format, but we can still analyze it as text/data
		result.Type = "Unknown"
	} else {
		result.Type = fileType
	}

	// Extract strings
	strings, err := extractStrings(path)
	if err == nil {
		result.Strings = strings
	}

	// Analyze based on file type
	switch result.Type {
	case "PE":
		analyzePE(path, result)
	case "ELF":
		analyzeELF(path, result)
	case "Mach-O":
		analyzeMachO(path, result)
	}

	return result, nil
}

// identifyFileType tries to determine the file type (PE, ELF, Mach-O, etc.)
func identifyFileType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the first few bytes to identify the file type
	header := make([]byte, 16)
	if _, err := file.Read(header); err != nil {
		return "", err
	}

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	// Check for PE files
	if bytes.Equal(header[0:2], []byte{0x4D, 0x5A}) { // "MZ"
		_, err := pe.NewFile(file)
		if err == nil {
			return "PE", nil
		}
	}

	// Check for ELF files
	if bytes.Equal(header[0:4], []byte{0x7F, 0x45, 0x4C, 0x46}) { // "\x7FELF"
		_, err := elf.NewFile(file)
		if err == nil {
			return "ELF", nil
		}
	}

	// Check for Mach-O files
	if bytes.Equal(header[0:4], []byte{0xFE, 0xED, 0xFA, 0xCE}) ||
		bytes.Equal(header[0:4], []byte{0xCE, 0xFA, 0xED, 0xFE}) ||
		bytes.Equal(header[0:4], []byte{0xFE, 0xED, 0xFA, 0xCF}) ||
		bytes.Equal(header[0:4], []byte{0xCF, 0xFA, 0xED, 0xFE}) {
		_, err := macho.NewFile(file)
		if err == nil {
			return "Mach-O", nil
		}
	}

	return "", errors.New("unknown file type")
}

// extractStrings extracts printable ASCII strings from a file
func extractStrings(path string) ([]string, error) {
	// Use the 'strings' command if available
	cmd := exec.Command("strings", path)
	output, err := cmd.Output()
	if err == nil {
		// Parse the output
		lines := strings.Split(string(output), "\n")
		var result []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if len(trimmed) >= 4 { // Only include strings of reasonable length
				result = append(result, trimmed)
			}
		}
		return result, nil
	}

	// Fallback: manual string extraction
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Find sequences of printable ASCII characters
	re := regexp.MustCompile(`[A-Za-z0-9/\-:.,_$%'()[\]<> ]{4,}`)
	matches := re.FindAllString(string(data), -1)

	return matches, nil
}

// analyzePE analyzes a PE (Windows) executable
func analyzePE(path string, info *FileInfo) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	peFile, err := pe.NewFile(file)
	if err != nil {
		return err
	}

	// Determine architecture
	switch peFile.Machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		info.Arch = "x86"
	case pe.IMAGE_FILE_MACHINE_AMD64:
		info.Arch = "x86_64"
	case pe.IMAGE_FILE_MACHINE_ARM:
		info.Arch = "ARM"
	case pe.IMAGE_FILE_MACHINE_ARMNT:
		info.Arch = "ARM Thumb-2"
	default:
		info.Arch = fmt.Sprintf("Unknown (%#x)", peFile.Machine)
	}

	// Extract imported functions
	imports, _ := peFile.ImportedSymbols()
	info.Imports = imports

	// Extract exported functions - simplified for compatibility
	// The Exports() method might not be available in all Go versions
	exportedSymbols := []string{}
	for _, exp := range peFile.Sections {
		if exp.Name == ".edata" || strings.Contains(exp.Name, "export") {
			// Found the export section
			data, err := exp.Data()
			if err == nil {
				// Simple extraction of potential function names
				potentialFuncs := regexp.MustCompile(`[A-Za-z0-9_]+`)
				matches := potentialFuncs.FindAllString(string(data), -1)
				for _, match := range matches {
					if len(match) > 3 && !strings.HasPrefix(match, "_") {
						exportedSymbols = append(exportedSymbols, match)
					}
				}
			}
		}
	}
	info.Exports = exportedSymbols

	// Extract sections
	for _, section := range peFile.Sections {
		isExec := section.Characteristics&pe.IMAGE_SCN_MEM_EXECUTE != 0
		info.Sections = append(info.Sections, Section{
			Name:         section.Name,
			Addr:         uint64(section.VirtualAddress),
			Size:         uint64(section.Size),
			Offset:       uint64(section.Offset),
			Flags:        fmt.Sprintf("%08x", section.Characteristics),
			IsExecutable: isExec,
		})
	}

	return nil
}

// analyzeELF analyzes an ELF (Linux) executable
func analyzeELF(path string, info *FileInfo) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	elfFile, err := elf.NewFile(file)
	if err != nil {
		return err
	}

	// Determine architecture
	switch elfFile.Machine {
	case elf.EM_386:
		info.Arch = "x86"
	case elf.EM_X86_64:
		info.Arch = "x86_64"
	case elf.EM_ARM:
		info.Arch = "ARM"
	case elf.EM_AARCH64:
		info.Arch = "ARM64"
	default:
		info.Arch = fmt.Sprintf("Unknown (%#x)", elfFile.Machine)
	}

	// Check if stripped
	syms, _ := elfFile.Symbols()
	info.IsStripped = len(syms) == 0

	// Extract imported functions from dynamic symbols
	dynsyms, _ := elfFile.DynamicSymbols()
	for _, sym := range dynsyms {
		if elf.ST_BIND(sym.Info) == elf.STB_GLOBAL && elf.ST_TYPE(sym.Info) == elf.STT_FUNC {
			info.Imports = append(info.Imports, sym.Name)
		}
	}

	// Extract sections
	for _, section := range elfFile.Sections {
		isExec := section.Flags&elf.SHF_EXECINSTR != 0
		info.Sections = append(info.Sections, Section{
			Name:         section.Name,
			Addr:         section.Addr,
			Size:         section.Size,
			Offset:       section.Offset,
			Flags:        fmt.Sprintf("%08x", section.Flags),
			IsExecutable: isExec,
		})
	}

	return nil
}

// analyzeMachO analyzes a Mach-O (macOS) executable
func analyzeMachO(path string, info *FileInfo) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	machoFile, err := macho.NewFile(file)
	if err != nil {
		return err
	}

	// Determine architecture based on CPU type constants
	// Constants may vary by Go version, using values directly for compatibility
	switch machoFile.Cpu {
	case 7: // CPU386 = 7
		info.Arch = "x86"
	case 0x01000007: // CPUAmd64 = 0x01000007
		info.Arch = "x86_64"
	case 12: // CPUArm = 12
		info.Arch = "ARM"
	case 0x01000012: // CPUArm64 = 0x01000012
		info.Arch = "ARM64"
	default:
		info.Arch = fmt.Sprintf("Unknown (%#x)", machoFile.Cpu)
	}

	// Extract imported functions
	imports, _ := machoFile.ImportedSymbols()
	info.Imports = imports

	// Extract sections
	for _, section := range machoFile.Sections {
		// S_ATTR_SOME_INSTRUCTIONS = 0x00000400
		isExec := section.Flags&0x00000400 != 0
		info.Sections = append(info.Sections, Section{
			Name:         section.Name,
			Addr:         section.Addr,
			Size:         section.Size,
			Offset:       uint64(section.Offset),
			Flags:        fmt.Sprintf("%08x", section.Flags),
			IsExecutable: isExec,
		})
	}

	return nil
}
