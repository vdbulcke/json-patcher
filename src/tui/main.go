package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vdbulcke/json-patcher/src/config"
	"github.com/vdbulcke/json-patcher/src/patcher"
)

type PatchItem struct {
	patch *config.Patch
}

// implement list item interface for UI
func (p PatchItem) Title() string {

	if strings.Contains(p.patch.Source, "/") {
		split := strings.Split(p.patch.Source, "/")
		return split[len(split)-1]
	}

	return p.patch.Source
}
func (p PatchItem) Description() string {
	return fmt.Sprintf("Source: %s\nDestination: %s", p.patch.Source, p.patch.Destination)
}
func (p PatchItem) FilterValue() string { return p.patch.Source }

// GeneratePatchItemList filter patches and format them to PatchItem
func GeneratePatchItemList(cfg *config.Config, opts *patcher.Options) ([]PatchItem, error) {
	var items []PatchItem

	for _, p := range cfg.Patches {
		// handle skips
		r, err := patcher.Skip(p, opts)
		if err != nil {
			return nil, err
		}

		if !r.Skip {
			pi := PatchItem{
				patch: p,
			}
			items = append(items, pi)

		}

	}
	return items, nil

}

func GenerateItemList(pitems []PatchItem) []list.Item {
	items := make([]list.Item, len(pitems))
	for i, p := range pitems {

		items[i] = list.Item(p)

	}

	return items
}

func StartUI(cfg *config.Config, opts *patcher.Options) error {

	pitems, err := GeneratePatchItemList(cfg, opts)
	if err != nil {
		return err
	}

	items := GenerateItemList(pitems)
	UIState = &State{
		List:                items,
		DisplayCurrentIndex: 0,
		ops:                 opts,
	}
	m, _ := InitList(UIState.List)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		// fmt.Println("Error running program:", err)
		// os.Exit(1)
		return err
	}

	return nil
}
