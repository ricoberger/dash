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
}

func (v *Variable) SetClient(client datasource.Client) {
	v.client = client
}

func (v *Variable) GetValues(variables map[string]string, start, end time.Time) ([]string, error) {
	query, err := datasource.QueryInterpolation(v.Query, variables)
	if err != nil {
		return nil, err
	}

	return v.client.GetVariableValues(query, v.Label, start, end)
}
