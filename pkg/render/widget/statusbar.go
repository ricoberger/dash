package widget

import (
	"fmt"
	"image"
	"strings"

	"github.com/ricoberger/dash/pkg/render/utils"

	ui "github.com/gizak/termui/v3"
)

type Statusbar struct {
	*ui.Block

	storage *utils.Storage
}

func NewStatusbar(termWidth, termHeight int, storage *utils.Storage) *Statusbar {
	bar := ui.NewBlock()
	bar.Border = false
	bar.PaddingTop = 5

	bar.SetRect(0, 0, termWidth, 1)

	return &Statusbar{
		bar,
		storage,
	}
}

func (s *Statusbar) Draw(buf *ui.Buffer) {
	// Render the background of the status bar. To render the whole background of the status bar in blue we have to
	// provide a string with the length of the status bar.
	buf.SetString(
		strings.Repeat(" ", s.Inner.Dx()),
		ui.NewStyle(ui.ColorBlack, ui.ColorBlue),
		image.Pt(s.Inner.Min.X, s.Inner.Min.Y),
	)

	// Left side of the statusbar contains the name of the current dashboard and the selected values for the specified
	// variables.
	dashboard := fmt.Sprintf("[D]ashboard: %s", s.storage.Dashboard().Name)
	buf.SetString(
		dashboard,
		ui.NewStyle(ui.ColorBlack, ui.ColorBlue),
		image.Pt(s.Inner.Min.X, s.Inner.Min.Y),
	)

	variables := fmt.Sprintf("[V]ariables: %s", strings.Join(s.storage.GetVariableValues(), ", "))
	buf.SetString(
		variables,
		ui.NewStyle(ui.ColorBlack, ui.ColorBlue),
		image.Pt(s.Inner.Min.X+len(dashboard)+2, s.Inner.Min.Y),
	)

	// Right side of the statusbar contains the selected time range and the refresh interval for new data.
	refreshInterval := fmt.Sprintf("[R]efresh: %s", s.storage.Refresh)
	buf.SetString(
		refreshInterval,
		ui.NewStyle(ui.ColorBlack, ui.ColorBlue),
		image.Pt(s.Inner.Max.X-len(refreshInterval), s.Inner.Min.Y),
	)

	interval := fmt.Sprintf("[I]nterval: %s", s.storage.Interval.Interval)
	buf.SetString(
		interval,
		ui.NewStyle(ui.ColorBlack, ui.ColorBlue),
		image.Pt(s.Inner.Max.X-len(refreshInterval)-2-len(interval), s.Inner.Min.Y),
	)
}
