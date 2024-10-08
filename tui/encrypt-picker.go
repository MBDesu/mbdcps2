package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type EncryptPicker struct{}

func (d EncryptPicker) New() (string, string) {
	zm := InitializeFilePicker("Choose a ROM .zip containing the key:", []string{".zip"})
	tzm, _ := tea.NewProgram(&zm, tea.WithAltScreen()).Run()
	zmm := tzm.(FilePickerModel)
	bm := InitializeFilePicker("Choose the decrypted .bin file:", []string{})
	tbm, _ := tea.NewProgram(&bm, tea.WithAltScreen()).Run()
	bmm := tbm.(FilePickerModel)
	return zmm.SelectedFile, bmm.SelectedFile
}
