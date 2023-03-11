package tui

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/vdbulcke/json-patcher/src/config"
	"github.com/vdbulcke/json-patcher/src/patcher"
	"gopkg.in/yaml.v2"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type PatchDisplay struct {
	viewport   viewport.Model
	wasApplied bool
	patch      *config.Patch
}

func newPatchDisplay(patch *config.Patch) (*PatchDisplay, error) {
	// const width = 78

	top, right, bottom, left := lipgloss.NewStyle().Margin(0, 2).GetMargin()
	vp := viewport.New(WindowSize.Width-left-right, WindowSize.Height-top-bottom-6)
	p := &PatchDisplay{
		viewport:   vp,
		patch:      patch,
		wasApplied: false,
	}
	p.viewport.Style = lipgloss.NewStyle().Align(lipgloss.Bottom)

	str, _ := glamour.Render(genMDPatchInfo(p.patch), "dark")
	p.viewport.SetContent(str)

	return p, nil
}
func (p PatchDisplay) Init() tea.Cmd {
	return nil
}

func (p PatchDisplay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		WindowSize = msg
		top, right, bottom, left := lipgloss.NewStyle().Margin(0, 2).GetMargin()
		p.viewport = viewport.New(WindowSize.Width-left-right, WindowSize.Height-top-bottom-6)

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c":
			return p, tea.Quit
		case "q", "esc":

			return InitList(UIState.List)
		case "a":
			if !p.wasApplied {

				// apply patch
				mdoc, err := patcher.Patch(p.patch, UIState.ops)
				if err != nil {
					// if error
					// generate markdown template
					data := genMDError(err)

					// render output
					str, _ := glamour.Render(data, "dark")
					p.viewport.SetContent(str)

				} else if p.patch.Destination == config.STDOUT {
					// generate markdown template
					data := genMDPreview(p.patch, mdoc, err)

					// render output
					str, _ := glamour.Render(data, "dark")
					p.viewport.SetContent(str)
				} else {

					// write patch
					err = patcher.WriteDestination(p.patch.Destination, mdoc, UIState.ops)

					// render output
					data := genMDApply(p.patch, err)
					// render output
					str, _ := glamour.Render(data, "dark")
					p.viewport.SetContent(str)

				}

				// update entry local state
				p.wasApplied = true
				// update state list removing applied entry
				UIState.List = remove(UIState.List, UIState.DisplayCurrentIndex)

			}

			var cmd tea.Cmd
			p.viewport, cmd = p.viewport.Update(msg)
			return p, cmd

		case "p":

			// apply patch
			mdoc, err := patcher.Patch(p.patch, UIState.ops)

			// generate markdown template
			data := genMDPreview(p.patch, mdoc, err)

			// render output
			str, _ := glamour.Render(data, "dark")
			p.viewport.SetContent(str)

			var cmd tea.Cmd
			p.viewport, cmd = p.viewport.Update(msg)
			return p, cmd
		case "t":
			// update option
			UIState.ops.AllowUnescapedHTML = !UIState.ops.AllowUnescapedHTML

			// re-instantiate a new display from updated options
			// from same patch
			newDisplay, _ := newPatchDisplay(p.patch)

			// set dummy "r" unbound key
			return newDisplay.Update("r")
		default:
			var cmd tea.Cmd
			p.viewport, cmd = p.viewport.Update(msg)
			return p, cmd
		}
	default:
		return p, nil
	}
	return p, nil
}

func (p PatchDisplay) View() string {
	return p.viewport.View() + p.helpView()
}

func (p PatchDisplay) helpView() string {

	if p.wasApplied {
		return helpStyle("\n  ↑/↓: Navigate • p: Preview •  q|esc: Back • ctrl+C: Quit\n")
	}
	return helpStyle("\n  ↑/↓: Navigate • p: Preview • t: Toggle unescaped HTML • a: Apply Patch • q|esc: Back • ctrl+C: Quit\n")
}
func remove(slice []list.Item, s int) []list.Item {
	return append(slice[:s], slice[s+1:]...)
}

func genMDPreview(p *config.Patch, pByte []byte, perr error) string {

	var out string

	// unmarshal and marshal json
	// for indent output
	var parsed json.RawMessage

	err := json.Unmarshal(pByte, &parsed)
	if err != nil {
		return genMDError(err)
	}

	// pretty print json
	pIndentByte, err := json.MarshalIndent(&parsed, "", "  ")
	if err != nil {
		return genMDError(err)
	}

	// revert HTML escape
	if UIState.ops.AllowUnescapedHTML {
		pIndentByte = bytes.Replace(pIndentByte, []byte("\\u0026"), []byte("&"), -1)
		pIndentByte = bytes.Replace(pIndentByte, []byte("\\u003c"), []byte("<"), -1)
		pIndentByte = bytes.Replace(pIndentByte, []byte("\\u003e"), []byte(">"), -1)
	}
	mdTemp := `
# Preview Patch

%s
%s
%s

### Info


| Summary 
| --- |  --- |
| Source | %s |
| Destination | %s |

## Options

- __Skip-Tags__: '%s'
- __allow-unescaped-html__: %t


> CurrentIndex %d 
`

	if perr != nil {
		return genMDError(perr)
	}
	out = fmt.Sprintf(mdTemp,
		"```json",
		string(pIndentByte),
		"```",
		p.Source,
		p.Destination,
		UIState.ops.SkipTags,
		UIState.ops.AllowUnescapedHTML,
		UIState.DisplayCurrentIndex,
	)

	return out

}

func genMDApply(p *config.Patch, applyErr error) string {

	var out string

	pByte, err := yaml.Marshal(&p)
	if err != nil {
		return genMDError(err)
	}
	mdTemp := `
# Apply Patch Result

%s
%s
%s


Result:  **Success**

## Options

- __Skip-Tags__: '%s'
- __allow-unescaped-html__: %t


> CurrentIndex %d 

`

	if applyErr != nil {
		return genMDError(applyErr)
	}

	out = fmt.Sprintf(mdTemp,
		"```yaml", string(pByte),
		"```",
		UIState.ops.SkipTags,
		UIState.ops.AllowUnescapedHTML,
		UIState.DisplayCurrentIndex,
	)
	return out

}

func genMDPatchInfo(p *config.Patch) string {

	pByte, err := yaml.Marshal(&p)
	if err != nil {
		return genMDError(err)
	}

	mdTemp := `
# Current Patch 

%s
%s
%s



## Options

- __Skip-Tags__: '%s'
- __allow-unescaped-html__: %t


> CurrentIndex %d 
`
	out := fmt.Sprintf(mdTemp,
		"```yaml", string(pByte),
		"```",
		UIState.ops.SkipTags,
		UIState.ops.AllowUnescapedHTML,
		UIState.DisplayCurrentIndex,
	)

	return out

}

func genMDError(err error) string {
	errTemplate := `
# Error 

%s
`
	return fmt.Sprintf(errTemplate, err.Error())
}
