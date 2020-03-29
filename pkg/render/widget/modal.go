package widget

import (
	"github.com/ricoberger/dash/pkg/render/utils"

	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
)

type ModalType string

const (
	ModalTypeDashboard ModalType = "Dashboard"
	ModalTypeVariable  ModalType = "Variable"
	ModalTypeInterval  ModalType = "Interval"
	ModalTypeRefresh   ModalType = "Refresh"
)

var intervals = []string{"5m", "15m", "30m", "1h", "3h", "6h", "12h", "24h", "2d", "7d", "30d"}
var refreshs = []string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}

type Modal struct {
	*w.List

	termWidth  int
	termHeight int
	storage    *utils.Storage
	options    *ModalOptions
}

type ModalOptions struct {
	Type          ModalType
	VariableIndex int
}

func NewModal(termWidth, termHeight int, storage *utils.Storage) *Modal {
	modal := w.NewList()
	modal.TextStyle = ui.NewStyle(ui.ColorYellow)
	modal.WrapText = false

	return &Modal{
		modal,

		termWidth,
		termHeight,
		storage,
		nil,
	}
}

func (m *Modal) SetDimensions(termWidth, termHeight int) {
	m.termWidth = termWidth
	m.termHeight = termHeight
}

func (m *Modal) Hide() {
	m.SetRect(0, 0, 0, 0)
}

func (m *Modal) Show(options *ModalOptions) bool {
	var index int
	m.options = options
	m.Title = string(m.options.Type)
	m.Rows = []string{}

	if m.options.Type == ModalTypeDashboard {
		index = m.storage.ActiveDashboard
		for _, dashboard := range m.storage.Dashboards {
			m.Rows = append(m.Rows, dashboard.Name)
		}
	} else if m.options.Type == ModalTypeVariable {
		m.options.VariableIndex = m.options.VariableIndex - 1
		if m.options.VariableIndex >= len(m.storage.Dashboard().Variables) {
			return false
		}

		variable := m.storage.Dashboard().Variables[m.options.VariableIndex]
		values, err := variable.GetValues(m.storage.VariableValues, m.storage.Interval.Start, m.storage.Interval.End)
		if err != nil {
			return false
		}

		for key, value := range values {
			m.Rows = append(m.Rows, value)

			if value == m.storage.VariableValues[variable.Name] {
				index = key
			}
		}
	} else if m.options.Type == ModalTypeInterval {
		m.Rows = intervals

		for key, value := range m.Rows {
			if value == m.storage.Interval.Interval {
				index = key
			}
		}
	} else if m.options.Type == ModalTypeRefresh {
		m.Rows = refreshs

		for key, value := range m.Rows {
			if value == m.storage.Refresh {
				index = key
			}
		}
	} else {
		return false
	}

	m.SelectedRow = index
	m.SetRect(m.termWidth/2-25, m.termHeight/2-10, m.termWidth/2+25, m.termHeight/2+10)
	return true
}

func (m *Modal) Select() ModalType {
	if m.options.Type == ModalTypeDashboard {
		m.storage.ChangeDashboard(m.SelectedRow)
	} else if m.options.Type == ModalTypeVariable {
		m.storage.ChangeVariable(m.storage.Dashboard().Variables[m.options.VariableIndex].Name, m.Rows[m.SelectedRow])
	} else if m.options.Type == ModalTypeInterval {
		m.storage.ChangeInterval(intervals[m.SelectedRow])
	} else if m.options.Type == ModalTypeRefresh {
		m.storage.ChangeRefresh(refreshs[m.SelectedRow])
	}

	return m.options.Type
}
