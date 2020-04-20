package widget

import (
	"fmt"
	"strings"

	"github.com/ricoberger/dash/pkg/render/utils"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

type Statusbar struct {
	*text.Text

	storage *utils.Storage
}

func NewStatusbar(termWidth int, storage *utils.Storage) (*Statusbar, error) {
	txt, err := text.New()
	if err != nil {
		return nil, err
	}

	statusbar := &Statusbar{txt, storage}
	statusbar.Update(termWidth)
	return statusbar, nil
}

func (s *Statusbar) Update(termWidth int) {
	s.Reset()

	var prefixedValues []string
	values := s.storage.GetVariableValues()

	for index, value := range values {
		prefixedValues = append(prefixedValues, fmt.Sprintf("[%d] %s", index+1, value))
	}

	dashboard := fmt.Sprintf(" [F1] Dashboard: %s", s.storage.Dashboard().Name)
	datasource := fmt.Sprintf(" [F2] Datasource: %s", s.storage.ActiveDatasource)
	variables := fmt.Sprintf(" [F3] Variables: %s", strings.Join(prefixedValues, ", "))
	interval := fmt.Sprintf(" [F4] Interval: %s", s.storage.Interval.Interval)
	refresh := fmt.Sprintf(" [F5] Refresh: %s ", s.storage.Refresh)
	spaces := strings.Repeat(" ", termWidth-len(dashboard)-len(datasource)-len(variables)-len(interval)-len(refresh))

	s.Write(dashboard+datasource+variables+spaces+interval+refresh, text.WriteCellOpts(cell.BgColor(cell.ColorBlue), cell.FgColor(cell.ColorBlack)))
}
