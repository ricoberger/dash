package dashboard

import (
	"time"

	"github.com/ricoberger/dash/pkg/datasource"
)

type Variable struct {
	client datasource.Client `yaml:"-"`

	Name       string `yaml:"name"`
	Datasource string `yaml:"datasource"`
	Query      string `yaml:"query"`
	Label      string `yaml:"label"`
	All        bool   `yaml:"all"`
}

func (v *Variable) SetClient(client datasource.Client) {
	v.client = client
}

func (v *Variable) GetValues(variables map[string]string, start, end time.Time) ([]string, error) {
	query, err := datasource.QueryInterpolation(v.Query, variables)
	if err != nil {
		return nil, err
	}

	values, err := v.client.GetVariableValues(query, v.Label, start, end)
	if err != nil {
		return nil, err
	}

	if v.All {
		return append([]string{".*"}, values...), nil
	}

	return values, nil
}
