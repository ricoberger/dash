package widget

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ricoberger/dash/pkg/render/utils"

	"github.com/mum4k/termdash/widgets/text"
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
	*text.Text

	storage *utils.Storage
	options *ModalOptions
	rows    []string
	index   string
}

type ModalOptions struct {
	Type          ModalType
	VariableIndex int
}

func NewModal(storage *utils.Storage) (*Modal, error) {
	modal, err := text.New()
	if err != nil {
		return nil, err
	}

	return &Modal{
		modal,

		storage,
		nil,
		nil,
		"",
	}, nil
}

func (m *Modal) show() bool {
	m.Reset()

	if m.options.Type == ModalTypeDashboard {
		for index, dashboard := range m.storage.Dashboards {
			m.rows = append(m.rows, fmt.Sprintf("%3d: %s", index, dashboard.Name))
		}
	} else if m.options.Type == ModalTypeVariable {
		if m.options.VariableIndex >= len(m.storage.Dashboard().Variables) {
			return false
		}

		variable := m.storage.Dashboard().Variables[m.options.VariableIndex]
		values, err := variable.GetValues(m.storage.VariableValues, m.storage.Interval.Start, m.storage.Interval.End)
		if err != nil {
			return false
		}

		for index, value := range values {
			m.rows = append(m.rows, fmt.Sprintf("%3d: %s", index, value))
		}
	} else if m.options.Type == ModalTypeInterval {
		for index, interval := range intervals {
			m.rows = append(m.rows, fmt.Sprintf("%3d: %s", index, interval))
		}
	} else if m.options.Type == ModalTypeRefresh {
		for index, refresh := range refreshs {
			m.rows = append(m.rows, fmt.Sprintf("%3d: %s", index, refresh))
		}
	} else {
		return false
	}

	if m.index == "" {
		err := m.Write(fmt.Sprintf("Selected index: \n\n%s", strings.Join(m.rows, "\n")))
		if err != nil {
			return false
		}
	} else {
		err := m.Write(fmt.Sprintf("Selected index: %s \n\n%s", m.index, strings.Join(m.rows, "\n")))
		if err != nil {
			return false
		}
	}

	return true
}

func (m *Modal) Show(options *ModalOptions) bool {
	m.options = options
	m.rows = nil
	m.index = ""
	return m.show()
}

func (m *Modal) SelectIndex(index string) bool {
	m.rows = nil
	m.index = m.index + index
	return m.show()
}

func (m *Modal) Select() (ModalType, error) {
	index, err := strconv.Atoi(m.index)
	if err != nil {
		return m.options.Type, err
	}

	if m.options.Type == ModalTypeDashboard {
		err := m.storage.ChangeDashboard(index)
		if err != nil {
			return m.options.Type, err
		}
	} else if m.options.Type == ModalTypeVariable {
		split := strings.Index(m.rows[index], ":")
		err := m.storage.ChangeVariable(m.storage.Dashboard().Variables[m.options.VariableIndex].Name, m.rows[index][split+2:])
		if err != nil {
			return m.options.Type, err
		}
	} else if m.options.Type == ModalTypeInterval {
		err := m.storage.ChangeInterval(intervals[index])
		if err != nil {
			return m.options.Type, err
		}
	} else if m.options.Type == ModalTypeRefresh {
		m.storage.ChangeRefresh(refreshs[index])
	}

	return m.options.Type, nil
}
