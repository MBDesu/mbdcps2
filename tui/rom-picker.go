package tui

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/MBDesu/mbdcps2/cps2rom"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (i item) filterValue() string { return "" }

type itemDelegate struct{}
type RomPicker struct{}

type RomPickerModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func initializeRomPicker() RomPickerModel {
	commands := []list.Item{}
	romDefs := cps2rom.RomDefinitions
	for k := range *romDefs {
		commands = append(commands, item{k})
	}
	slices.SortFunc(commands, func(a list.Item, b list.Item) int {
		return strings.Compare(a.FilterValue(), b.FilterValue())
	})
	l := list.New(commands, ItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Which ROM set are you working with?"
	l.SetShowStatusBar(false)
	l.Styles.PaginationStyle = PaginationStyle
	l.Styles.HelpStyle = HelpStyle
	return RomPickerModel{list: l}
}

func (m RomPicker) New() string {
	w := initializeRomPicker()
	o, err := tea.NewProgram(w, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return o.(RomPickerModel).choice
}

func (m RomPickerModel) Init() tea.Cmd {
	return nil
}

func (m RomPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m RomPickerModel) View() string {
	if m.choice != "" {
		return fmt.Sprintf("%s? Sounds good to me.", m.choice)
	}
	if m.quitting {
		return fmt.Sprintf("Bye")
	}
	return "\n" + m.list.View()
}
