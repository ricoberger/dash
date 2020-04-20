package datasource

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	fLog "github.com/ricoberger/dash/pkg/log"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Prometheus struct {
	v1api   v1.API
	options Options
}

type basicAuthTransport struct {
	Transport http.RoundTripper

	username string
	password string
}

type tokenAuthTransporter struct {
	Transport http.RoundTripper

	token string
}

func (bat basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(bat.username, bat.password)
	return bat.Transport.RoundTrip(req)
}

func (tat tokenAuthTransporter) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+tat.token)
	return tat.Transport.RoundTrip(req)
}

func NewPrometheusClient(datasource Datasource) (*Prometheus, error) {
	roundTripper := api.DefaultRoundTripper

	if datasource.Auth.Username != "" && datasource.Auth.Password != "" {
		roundTripper = basicAuthTransport{
			Transport: roundTripper,
			username:  datasource.Auth.Username,
			password:  datasource.Auth.Password,
		}
	}

	if datasource.Auth.Token != "" {
		roundTripper = tokenAuthTransporter{
			Transport: roundTripper,
			token:     datasource.Auth.Token,
		}
	}

	client, err := api.NewClient(api.Config{
		Address:      datasource.URL,
		RoundTripper: roundTripper,
	})

	if err != nil {
		return nil, err
	}

	return &Prometheus{
		v1api:   v1.NewAPI(client),
		options: datasource.Options,
	}, nil
}

func (p *Prometheus) GetVariableValues(query, label string, start, end time.Time) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	labelSets, _, err := p.v1api.Series(ctx, []string{query}, start, end)
	if err != nil {
		return nil, err
	}

	var values []string

	for _, labelSet := range labelSets {
		if value, ok := labelSet[model.LabelName(label)]; ok {
			values = appendIfMissing(values, string(value))
		}
	}

	return values, nil
}

func (p *Prometheus) GetData(queries, labels []string, start, end time.Time) (*Data, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var series []Series
	var timestamps map[int]string
	timestamps = make(map[int]string)

	timeRange := getTimeRange(p.options, start, end)

	for i, query := range queries {
		result, _, err := p.v1api.QueryRange(ctx, query, timeRange)
		if err != nil {
			return nil, err
		}

		data, ok := result.(model.Matrix)
		if !ok {
			return nil, fmt.Errorf("unsupported result format: %s", result.Type().String())
		}

		for j, d := range data {
			fLog.Debugf("query %s returned %d points and the following labels %v", query, len(d.Values), d.Metric)

			var points []float64
			var returnedLabels map[string]string
			returnedLabels = make(map[string]string)

			for key, value := range d.Metric {
				returnedLabels[string(key)] = string(value)
			}

			for key, value := range d.Values {
				if i == 0 && j == 0 {
					timestamps[key] = value.Timestamp.Time().Format("01/02 15:04")
				}
				points = append(points, float64(value.Value))
			}

			series = append(series, Series{
				Label:  getLabel(labels[i], returnedLabels),
				Points: points,
			})
		}
	}

	return &Data{
		Timestamps: timestamps,
		Series:     series,
	}, nil
}

func (p *Prometheus) GetTableData(queries, labels []string) (*TableData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var tableData TableData
	tableData = make(map[string]map[string]interface{})

	now := time.Now()

	for i, query := range queries {
		result, _, err := p.v1api.Query(ctx, query, now)
		if err != nil {
			return nil, err
		}

		data, ok := result.(model.Vector)
		if !ok {
			return nil, fmt.Errorf("unsupported result format: %s", result.Type().String())
		}

		for _, d := range data {
			var returnedLabels map[string]string
			returnedLabels = make(map[string]string)

			for key, value := range d.Metric {
				returnedLabels[string(key)] = string(value)
			}

			joinValue := getLabel(labels[i], returnedLabels)
			for key, value := range returnedLabels {
				if _, ok := tableData[joinValue]; !ok {
					tableData[joinValue] = make(map[string]interface{})
				}

				if _, ok := tableData[joinValue][key]; !ok {
					tableData[joinValue][key] = value
				}

				tableData[joinValue][fmt.Sprintf("value_%d", i)] = float64(d.Value)
			}
		}
	}

	return &tableData, nil
}

func (p *Prometheus) GetSuggestions() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, _, err := p.v1api.LabelValues(ctx, "__name__")
	if err != nil {
		return nil, err
	}

	var values []string
	for _, val := range result {
		values = append(values, string(val))
	}

	return values, nil
}

func getTimeRange(options Options, start, end time.Time) v1.Range {
	var step = 10 * time.Second
	if options.MaxPoints != 0 {
		step = time.Duration((end.Unix()-start.Unix())/options.MaxPoints) * time.Second
	} else if options.Step != 0 {
		step = time.Duration(options.Step) * time.Second
	}

	return v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}
}

func getLabel(label string, labels map[string]string) string {
	value, err := QueryInterpolation(label, labels)
	if err != nil || label == "" {
		var values []string
		for key, value := range labels {
			values = append(values, key+"="+value)
		}
		return strings.Join(values, ", ")
	}

	return value
}

func appendIfMissing(items []string, item string) []string {
	for _, ele := range items {
		if ele == item {
			return items
		}
	}

	return append(items, item)
}
