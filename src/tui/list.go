package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vdbulcke/json-patcher/src/patcher"
)

var (
	UIState *State

	docStyle = lipgloss.NewStyle().Margin(1, 2)
	// WindowSize store the size of the terminal window
	WindowSize tea.WindowSizeMsg
)

type State struct {
	// list of item from list view
	List []list.Item
	// current index from the list for
	// viewport view
	DisplayCurrentIndex int

	ops *patcher.Options
}

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
	viewItem         key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		// insertItem: key.NewBinding(
		// 	key.WithKeys("a"),
		// 	key.WithHelp("a", "add item"),
		// ),
		// viewItem: key.NewBinding(
		// 	key.WithKeys("v"),
		// 	key.WithHelp("v", "view item"),
		// ),

		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	quitting     bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		WindowSize = msg
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.delegateKeys.choose):
			p, ok := m.list.SelectedItem().(PatchItem)
			if ok {
				// set current index in state
				UIState.DisplayCurrentIndex = m.list.Index()

				entry, err := newPatchDisplay(p.patch)
				if err != nil {
					m.quitting = true
					return m, tea.Quit
				}

				return entry.Update(msg)

			}
		}

		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {

	return docStyle.Render(m.list.View())
}

func InitList(items []list.Item) (tea.Model, tea.Cmd) {

	listKeys := newListKeyMap()
	delegateKeys := newDelegateKeyMap()
	keyDelegate := newItemDelegate(delegateKeys)
	m := model{
		list:         list.New(items, keyDelegate, 8, 8),
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
	m.list.Title = "List of Valid Patches"
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.viewItem,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	if WindowSize.Height != 0 {
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(WindowSize.Width-left-right, WindowSize.Height-top-bottom-1)
	}

	return m, nil
}
