//go:build headless

package main

import "errors"

func (a *App) SaveAsLines(lines []Line) error {
	return errors.New("native save dialog is unavailable in headless builds")
}

func (a *App) OpenIni() (map[string]map[string]string, error) {
	return nil, errors.New("native open dialog is unavailable in headless builds")
}

func (a *App) SaveAs(data map[string]map[string]string) error {
	return errors.New("native save dialog is unavailable in headless builds")
}
