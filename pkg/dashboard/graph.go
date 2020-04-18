package dashboard

import (
	"time"

	"github.com/ricoberger/dash/pkg/datasource"
)

type Graph struct {
	Width      int     `yaml:"width"`
	Datasource string  `yaml:"datasource"`
	Type       string  `yaml:"type"`
	Title      string  `yaml:"title"`
	Queries    []Query `yaml:"queries"`
	Options    Options `yaml:"options"`
}

type Query struct {
	Query string `yaml:"query"`
	Label string `yaml:"label"`
}

type Options struct {
	Unit       string            `yaml:"unit"`
	Stats      []string          `yaml:"stats"`
	Decimals   int               `yaml:"decimals"`
	Thresholds []float64         `yaml:"thresholds"`
	Colors     []string          `yaml:"colors"`
	Legend     string            `yaml:"legend"`
	Mappings   map[string]string `yaml:"mappings"`
	Columns    []Column          `yaml:"columns"`
}

type Column struct {
	Name   string `yaml:"name"`
	Header string `yaml:"header"`
}

func (g *Graph) GetData(ds datasource.Client, variables map[string]string, start, end time.Time) (*datasource.Data, error) {
	var queries []string
	var labels []string

	for _, query := range g.Queries {
		q, err := datasource.QueryInterpolation(query.Query, variables)
		if err != nil {
			return nil, err
		}

		queries = append(queries, q)
		labels = append(labels, query.Label)
	}

	return ds.GetData(queries, labels, start, end)
}

func (g *Graph) GetTableData(ds datasource.Client, variables map[string]string) (*datasource.TableData, error) {
	var queries []string
	var labels []string

	for _, query := range g.Queries {
		q, err := datasource.QueryInterpolation(query.Query, variables)
		if err != nil {
			return nil, err
		}

		queries = append(queries, q)
		labels = append(labels, query.Label)
	}

	return ds.GetTableData(queries, labels)
}
