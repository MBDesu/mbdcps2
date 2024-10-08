package tui

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type FilePickerModel struct {
	filepicker   filepicker.Model
	SelectedFile string
	title        string
	quitting     bool
	err          error
}

type Filepicker struct{}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func InitializeFilePicker(title string, allowedFiletypes []string) FilePickerModel {
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fp := filepicker.New()
	fp.AllowedTypes = allowedFiletypes
	fp.CurrentDirectory = ex

	return FilePickerModel{
		filepicker: fp,
		title:      title,
	}
}

func (m Filepicker) New(title string, allowedFiletypes []string) string {
	fp := InitializeFilePicker(title, allowedFiletypes)
	fpm, _ := tea.NewProgram(&fp, tea.WithAltScreen()).Run()
	fpmm := fpm.(FilePickerModel)
	return fpmm.SelectedFile
}

func (m FilePickerModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m FilePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)
	cmds = append(cmds, cmd)

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.SelectedFile = path
		return m, tea.Quit
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid.")
		m.SelectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, tea.Batch(cmds...)
}

func (m FilePickerModel) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.SelectedFile == "" {
		s.WriteString(list.DefaultStyles().Title.Render(m.title))
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.SelectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}
