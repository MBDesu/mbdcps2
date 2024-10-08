package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type DecryptPicker struct{}

func (d DecryptPicker) New() string {
	m := InitializeFilePicker("Choose a ROM .zip to decrypt:", []string{".zip"})
	tm, _ := tea.NewProgram(&m, tea.WithAltScreen()).Run()
	mm := tm.(FilePickerModel)
	return mm.SelectedFile
}
