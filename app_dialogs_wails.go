//go:build !headless

package main

import (
	"errors"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SaveAsLines shows a native save dialog and writes the provided lines to the
// selected path. The internal state is updated if the user selects a path.
func (a *App) SaveAsLines(lines []Line) error {
	if lines == nil {
		return errors.New("no lines provided to save")
	}
	if a.ctx == nil {
		return errors.New("application not ready")
	}
	defaultName := "config.ini"
	if a.currentFilePath != "" {
		defaultName = filepath.Base(a.currentFilePath)
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "儲存為",
		DefaultFilename: defaultName,
		Filters:         []runtime.FileFilter{{DisplayName: "INI Files (*.ini)", Pattern: "*.ini"}},
	})
	if err != nil || path == "" {
		return err
	}
	newData, err := a.helperSaveLines(lines, path)
	if err != nil {
		return err
	}
	a.currentFilePath = path
	a.lines = lines
	a.data = newData
	return nil
}

// OpenIni shows a native file picker allowing the user to select an INI file.
func (a *App) OpenIni() (map[string]map[string]string, error) {
	if a.ctx == nil {
		return nil, errors.New("application not ready")
	}
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "選擇 INI 檔案",
		Filters: []runtime.FileFilter{{DisplayName: "INI Files (*.ini)", Pattern: "*.ini"}},
	})
	if err != nil || selection == "" {
		return nil, err
	}
	return a.LoadIni(selection)
}

// SaveAs shows a native save file dialog to let the user choose where to
// write the provided INI data.
func (a *App) SaveAs(data map[string]map[string]string) error {
	if data == nil {
		return errors.New("no data provided to save")
	}
	if a.ctx == nil {
		return errors.New("application not ready")
	}
	defaultName := "config.ini"
	if a.currentFilePath != "" {
		defaultName = filepath.Base(a.currentFilePath)
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "儲存為",
		DefaultFilename: defaultName,
		Filters:         []runtime.FileFilter{{DisplayName: "INI Files (*.ini)", Pattern: "*.ini"}},
	})
	if err != nil || path == "" {
		return err
	}
	if err := a.SaveIni(path, data); err != nil {
		return err
	}
	a.currentFilePath = path
	return nil
}
