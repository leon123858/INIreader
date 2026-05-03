package main

import (
	"bufio"
	"context"
	"errors"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

// LineType represents the type of an INI line.
// It can be one of:
//
//	"blank"    - an empty or whitespace-only line
//	"comment"  - a line that starts with a comment symbol (';' or '#')
//	"section"  - a section header like [section]
//	"key"      - a key/value pair, possibly commented out
//
// Additional types may be added in the future.
//
// When saving, the line type along with the Commented flag determines how
// the line is reconstructed.
const (
	LineTypeBlank   = "blank"
	LineTypeComment = "comment"
	LineTypeSection = "section"
	LineTypeKey     = "key"
)

// Line represents a single line in an INI file. It retains the original
// position and allows editing of the section name, key, value, and comment
// status. Commented indicates whether the line is commented out (i.e. begins
// with ';' or '#') when saved. For comment-only lines, the Value field
// holds the comment text. For section headers, Section holds the name and
// Key/Value are empty. For key/value pairs, Section, Key and Value are
// populated and the line type is "key".
type Line struct {
	Section   string `json:"section"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Commented bool   `json:"commented"`
	Type      string `json:"type"`
}

// App implements the backend logic for the INI reader application.
// It exposes methods that can be called from the frontend via Wails.  The
// application maintains the currently loaded file path and a simple
// representation of the INI data as a nested map.  Clients can load an INI
// file, retrieve its contents, modify the data and write it back to disk.
//
// Methods returning an error will automatically propagate the error back to
// JavaScript as a rejected promise.  Primitive types and Go maps will be
// marshalled to JSON.
type App struct {
	ctx             context.Context
	initialFilePath string
	currentFilePath string
	// data holds the parsed INI data in a simple nested map form. It is kept
	// for backward compatibility with previous versions of the UI that
	// operated on a map-of-map representation. Newer functionality that
	// preserves ordering and comments uses the lines field.
	data map[string]map[string]string
	// lines represents the contents of the loaded INI file as an ordered
	// slice of Line structs. Each entry corresponds to a line in the file
	// and retains information about the section, key, value, whether the
	// line is commented, and the line type (section header, key/value pair,
	// comment-only line or blank). This allows the frontend to display and
	// manipulate commented lines and preserve ordering when saving.
	lines []Line
}

// helperSaveLines rebuilds an INI configuration from the provided lines and
// writes it to the given path. It returns the reconstructed data map and
// any error encountered. The function encapsulates all logic for
// reconstructing sections, keys and associated comments using the
// go-ini/ini package. When using this helper, callers should update
// internal state (such as a.currentFilePath) themselves upon success.
func (a *App) helperSaveLines(lines []Line, path string) (map[string]map[string]string, error) {
	if lines == nil {
		return nil, errors.New("no lines provided to save")
	}
	var sb strings.Builder
	// Iterate over each line and reconstruct the literal representation.
	for _, l := range lines {
		switch l.Type {
		case LineTypeBlank:
			// Write an empty line
			sb.WriteString("\n")
		case LineTypeComment:
			// Write pure comment lines. Always prefix with ';'.
			if l.Value != "" {
				sb.WriteString(";" + l.Value + "\n")
			} else {
				sb.WriteString(";\n")
			}
		case LineTypeSection:
			// Section header. If commented, prefix with ';'.
			line := "[" + l.Section + "]"
			if l.Commented {
				sb.WriteString(";" + line + "\n")
			} else {
				sb.WriteString(line + "\n")
			}
		case LineTypeKey:
			// Key/value pair. Build the expression. An empty key indicates a
			// malformed line; write only the value.
			var expr string
			if l.Key != "" {
				expr = l.Key + " = " + l.Value
			} else {
				expr = l.Value
			}
			if l.Commented {
				sb.WriteString(";" + expr + "\n")
			} else {
				sb.WriteString(expr + "\n")
			}
		default:
			// Unknown types treated as comment lines
			if l.Value != "" {
				sb.WriteString(";" + l.Value + "\n")
			} else {
				sb.WriteString(";\n")
			}
		}
	}
	// Write the constructed content to the file
	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		return nil, err
	}
	// Parse the file again to build a simple nested map of uncommented keys
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	newData := make(map[string]map[string]string)
	for _, secName := range cfg.SectionStrings() {
		sec := cfg.Section(secName)
		kv := make(map[string]string)
		for _, keyName := range sec.KeyStrings() {
			kv[keyName] = sec.Key(keyName).String()
		}
		newData[secName] = kv
	}
	return newData, nil
}

// NewApp creates a new App instance.  It accepts an optional
// initial file path which, if provided, is loaded during startup.  The
// initial file path is typically set by main.go after parsing command line
// arguments passed via a file association on Windows.
func NewApp(initialFile string) *App {
	return &App{
		initialFilePath: initialFile,
		data:            make(map[string]map[string]string),
		lines:           make([]Line, 0),
	}
}

// startup is called by Wails when the application is ready.  It stores
// the context for later use and loads the initial file if one was specified.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.initialFilePath != "" {
		// Ignoring the returned data here because the frontend will call GetData()
		_, _ = a.LoadIni(a.initialFilePath)
	}
}

// GetData returns the current INI data as a nested map.  The
// outer keys are section names and the inner keys are key names.
func (a *App) GetData() map[string]map[string]string {
	// If no data has been loaded yet, return an empty map rather than nil
	if a.data == nil {
		a.data = make(map[string]map[string]string)
	}
	return a.data
}

// GetLines returns the current INI lines preserving order and comment status.
// Each element in the returned slice corresponds to a line in the INI file.
// Lines are returned even if they are commented, blank or malformed.  This
// method can be called from the frontend to populate a table that allows
// editing and toggling comments.
func (a *App) GetLines() []Line {
	if a.lines == nil {
		a.lines = make([]Line, 0)
	}
	return a.lines
}

// SaveLines writes the provided lines back to the current file.  It updates
// the internal line slice and rebuilds the simple map data representation.
// If no file is currently loaded, an error is returned.  When reconstructing
// the file, the comment status and line type determine how each line is
// written.  Comment symbols are semicolons (;) for commented lines.
func (a *App) SaveLines(lines []Line) error {
	// Ensure there is a file loaded
	if a.currentFilePath == "" {
		return errors.New("no file loaded")
	}
	// Use helper to save lines to the existing path
	data, err := a.helperSaveLines(lines, a.currentFilePath)
	if err != nil {
		return err
	}
	a.lines = lines
	a.data = data
	return nil
}

// CurrentFile returns the currently loaded file path.  If no file has been
// loaded yet, an empty string is returned.
func (a *App) CurrentFile() string {
	return a.currentFilePath
}

// LoadIni reads the specified INI file and populates the application's
// internal data structure.  It returns the loaded data on success.  If the
// file cannot be read or parsed, an error is returned.  Note that this
// method does not store the returned data in the App unless there is no
// error.  The loaded file becomes the current file.
func (a *App) LoadIni(path string) (map[string]map[string]string, error) {
	// Read file manually to preserve ordering and comments
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var lines []Line
	currentSection := ""
	// Temporary map to build a simple key/value representation (excluding commented keys)
	result := make(map[string]map[string]string)
	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		// Detect blank line
		if trimmed == "" {
			lines = append(lines, Line{Type: LineTypeBlank, Commented: false})
			continue
		}
		// Detect commented line (starting with ';' or '#')
		if strings.HasPrefix(trimmed, ";") || strings.HasPrefix(trimmed, "#") {
			// Remove the comment symbol and trim spaces
			content := strings.TrimSpace(trimmed[1:])
			// Determine if this commented line contains a section or key/value
			if strings.HasPrefix(content, "[") && strings.Contains(content, "]") {
				// Commented section header: mark as section with Commented=true
				secName := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(content, "["), "]"))
				lines = append(lines, Line{Type: LineTypeSection, Section: secName, Commented: true})
			} else if strings.Contains(content, "=") {
				// Commented key/value pair: mark as key with Commented=true
				parts := strings.SplitN(content, "=", 2)
				key := strings.TrimSpace(parts[0])
				value := ""
				if len(parts) > 1 {
					value = strings.TrimSpace(parts[1])
				}
				lines = append(lines, Line{Type: LineTypeKey, Section: currentSection, Key: key, Value: value, Commented: true})
			} else {
				// Pure comment line: keep it as a comment line.  Use Commented=true to
				// indicate this line is always a comment and cannot be toggled to
				// uncommented.  The Value field holds the comment text for display.
				lines = append(lines, Line{Type: LineTypeComment, Value: content, Commented: true})
			}
			continue
		}
		// Detect section header
		if strings.HasPrefix(trimmed, "[") && strings.Contains(trimmed, "]") {
			endIdx := strings.Index(trimmed, "]")
			name := strings.TrimSpace(trimmed[1:endIdx])
			lines = append(lines, Line{Type: LineTypeSection, Section: name, Commented: false})
			currentSection = name
			continue
		}
		// Detect key/value pair
		if strings.Contains(trimmed, "=") {
			parts := strings.SplitN(trimmed, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := ""
			if len(parts) > 1 {
				value = strings.TrimSpace(parts[1])
			}
			lines = append(lines, Line{Type: LineTypeKey, Section: currentSection, Key: key, Value: value, Commented: false})
			// Populate result map with uncommented keys only
			if key != "" {
				if _, ok := result[currentSection]; !ok {
					result[currentSection] = make(map[string]string)
				}
				result[currentSection][key] = value
			}
			continue
		}
		// Any other format is treated as a plain comment line
		lines = append(lines, Line{Type: LineTypeComment, Value: trimmed, Commented: false})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	a.lines = lines
	a.data = result
	a.currentFilePath = path
	return result, nil
}

// SaveIni writes the provided data to the given file path.  If the path
// matches the currently loaded file, the application's internal state is
// updated as well.  When saving, sections are written in the order they
// appear in the map iteration.  If the file cannot be written, an error
// is returned.  Passing a nil or empty data map results in an error.
func (a *App) SaveIni(path string, data map[string]map[string]string) error {
	if data == nil {
		return errors.New("no data provided to save")
	}
	cfg := ini.Empty()
	for section, kv := range data {
		sec, _ := cfg.NewSection(section)
		for k, v := range kv {
			sec.NewKey(k, v)
		}
	}
	if err := cfg.SaveTo(path); err != nil {
		return err
	}
	// Update current state if this was the current file
	if a.currentFilePath == path {
		a.data = data
	}
	return nil
}
