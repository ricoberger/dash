package utils

import (
	"time"

	"github.com/ricoberger/dash/pkg/dashboard"
	"github.com/ricoberger/dash/pkg/datasource"
	fLog "github.com/ricoberger/dash/pkg/log"
)

const (
	initialActiveDashboard = 0
)

type Interval struct {
	Interval string
	Start    time.Time
	End      time.Time
}

type Storage struct {
	Datasources      map[string]datasource.Client
	Dashboards       []dashboard.Dashboard
	ActiveDatasource string
	ActiveDashboard  int
	Interval         Interval
	Refresh          string
	VariableValues   map[string]string
}

func (s *Storage) loadVariables() error {
	for _, variable := range s.Dashboards[s.ActiveDashboard].Variables {
		values, err := variable.GetValues(s.Datasource(), s.VariableValues, s.Interval.Start, s.Interval.End)
		if err != nil {
			return err
		}

		if len(values) == 0 {
			s.VariableValues[variable.Name] = ""
			continue
		} else {
			if value, ok := s.VariableValues[variable.Name]; !ok {
				s.VariableValues[variable.Name] = values[0]
			} else {
				if !valueExists(value, values) {
					s.VariableValues[variable.Name] = values[0]
				}
			}
		}
	}

	return nil
}

func (s *Storage) Datasource() datasource.Client {
	return s.Datasources[s.ActiveDatasource]
}

func (s *Storage) Dashboard() dashboard.Dashboard {
	return s.Dashboards[s.ActiveDashboard]
}

func (s *Storage) ChangeDatasource(active string) error {
	fLog.Debugf("change datasource to %s", active)
	s.ActiveDatasource = active
	s.VariableValues = make(map[string]string)
	return s.loadVariables()
}

func (s *Storage) ChangeDashboard(active int) error {
	fLog.Debugf("change dashboard index to %d", active)
	s.ActiveDashboard = active

	if _, ok := s.Datasources[s.Dashboards[active].DefaultDatasource]; ok {
		s.ActiveDatasource = s.Dashboards[active].DefaultDatasource
	}

	s.VariableValues = make(map[string]string)
	return s.loadVariables()
}

func (s *Storage) GetVariableValues() []string {
	var values []string

	for _, variable := range s.Dashboards[s.ActiveDashboard].Variables {
		if value, ok := s.VariableValues[variable.Name]; ok {
			values = append(values, value)
		}
	}

	return values
}

func (s *Storage) ChangeVariable(name, value string) error {
	fLog.Debugf("change variable %s to %s", name, value)
	s.VariableValues[name] = value
	return s.loadVariables()
}

func (s *Storage) ChangeInterval(interval string) error {
	fLog.Debugf("change interval to %s", interval)
	start, end := GetStartAndEndTime(interval)
	s.Interval.Interval = interval
	s.Interval.Start = start
	s.Interval.End = end
	return s.loadVariables()
}

func (s *Storage) GetRefresh() time.Duration {
	switch s.Refresh {
	case "5s":
		return 5 * time.Second
	case "10s":
		return 10 * time.Second
	case "30s":
		return 30 * time.Second
	case "1m":
		return 1 * time.Minute
	case "5m":
		return 5 * time.Minute
	case "15m":
		return 15 * time.Minute
	case "30m":
		return 30 * time.Minute
	case "1h":
		return 1 * time.Hour
	case "2h":
		return 2 * time.Hour
	case "1d":
		return 24 * time.Hour
	default:
		return 5 * time.Minute
	}
}

func (s *Storage) ChangeRefresh(refresh string) {
	fLog.Debugf("change refresh to %s", refresh)
	s.Refresh = refresh
}

func (s *Storage) RefreshInterval() {
	start, end := GetStartAndEndTime(s.Interval.Interval)
	s.Interval.Start = start
	s.Interval.End = end
}

func NewStorage(datasources map[string]datasource.Client, dashboards []dashboard.Dashboard, initialInterval, initialRefresh string) (*Storage, error) {
	start, end := GetStartAndEndTime(initialInterval)

	var initialActiveDatasource string
	if _, ok := datasources[dashboards[initialActiveDashboard].DefaultDatasource]; ok {
		initialActiveDatasource = dashboards[initialActiveDashboard].DefaultDatasource
	} else {
		for key := range datasources {
			initialActiveDatasource = key
			break
		}
	}

	s := &Storage{
		Datasources:      datasources,
		Dashboards:       dashboards,
		ActiveDatasource: initialActiveDatasource,
		ActiveDashboard:  initialActiveDashboard,
		Interval: Interval{
			Interval: initialInterval,
			Start:    start,
			End:      end,
		},
		Refresh:        initialRefresh,
		VariableValues: make(map[string]string),
	}

	s.loadVariables()

	return s, nil
}

func valueExists(value string, values []string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}
