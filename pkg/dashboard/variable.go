package dashboard

import (
	"time"

	"github.com/ricoberger/dash/pkg/datasource"
)

type Variable struct {
	Name  string `yaml:"name"`
	Query string `yaml:"query"`
	Label string `yaml:"label"`
	All   bool   `yaml:"all"`
}

func (v *Variable) GetValues(ds datasource.Client, variables map[string]string, start, end time.Time) ([]string, error) {
	query, err := datasource.QueryInterpolation(v.Query, variables)
	if err != nil {
		return nil, err
	}

	values, err := ds.GetVariableValues(query, v.Label, start, end)
	if err != nil {
		return nil, err
	}

	if v.All {
		return append([]string{".*"}, values...), nil
	}

	return values, nil
}
