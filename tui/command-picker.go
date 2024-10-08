package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 24
const defaultWidth = 20

var (
	TitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	PaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	QuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	title string
}

func (i item) FilterValue() string { return i.title }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(item)
	i := ""
	if !ok {
		return
	} else {
		i = item.title
	}

	str := fmt.Sprintf("%d. %s", index+1, i)
	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprintf(w, fn(str))
}

type CommandPickerModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func InitializeCommandPicker() CommandPickerModel {
	commands := []list.Item{
		item{"Decrypt"},
		item{"Encrypt"},
		item{"Patch"},
		item{"Diff"},
	}
	l := list.New(commands, ItemDelegate{}, defaultWidth, listHeight)
	l.Title = "What would you like to do?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = PaginationStyle
	l.Styles.HelpStyle = HelpStyle
	return CommandPickerModel{list: l}
}

func (m CommandPickerModel) Init() tea.Cmd {
	return nil
}

func (m CommandPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m CommandPickerModel) View() string {
	if m.choice != "" {
		return ""
	}
	if m.quitting {
		return ""
	}
	return "\n" + m.list.View()
}
