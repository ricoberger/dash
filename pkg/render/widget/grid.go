package widget

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ricoberger/dash/pkg/dashboard"
	"github.com/ricoberger/dash/pkg/datasource"
	"github.com/ricoberger/dash/pkg/render/utils"
	cPlot "github.com/ricoberger/dash/pkg/render/widget/plot"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Grid struct {
	*ui.Grid

	storage *utils.Storage
}

func NewGrid(termWidth, termHeight int, storage *utils.Storage) *Grid {
	grid := ui.NewGrid()
	grid.SetRect(0, 1, termWidth, termHeight)

	var rows []interface{}

	for _, row := range storage.Dashboard().Rows {
		var cols []interface{}

		for _, graph := range row.Graphs {
			var components []interface{}

			data, err := graph.GetData(storage.VariableValues, storage.Interval.Start, storage.Interval.End)
			if err != nil {
				p := widgets.NewParagraph()
				p.Title = graph.Title
				p.Text = fmt.Sprintf("Could not load data: %s", err.Error())
				components = append(components, p)
			} else {
				switch graph.Type {
				case "singlestat":
					components = append(components, singlestat(graph, data)...)
				case "gauge":
					components = append(components, gauge(graph, data)...)
				case "sparkline":
					components = append(components, sparkline(graph, data)...)
				case "plot":
					components = append(components, plot(graph, data)...)
				}
			}

			cols = append(cols, ui.NewCol(graph.Width, components...))
		}

		rows = append(rows, ui.NewRow(row.Height, cols...))
	}

	grid.Set(rows...)

	return &Grid{
		grid,

		storage,
	}
}

func (g *Grid) Refresh() {
	g.storage.RefreshInterval()
	var rows []interface{}

	for _, row := range g.storage.Dashboard().Rows {
		var cols []interface{}

		for _, graph := range row.Graphs {
			var components []interface{}

			data, err := graph.GetData(g.storage.VariableValues, g.storage.Interval.Start, g.storage.Interval.End)
			if err != nil {
				p := widgets.NewParagraph()
				p.Title = graph.Title
				p.Text = fmt.Sprintf("Could not load data: %s", err.Error())
				components = append(components, p)
			} else {
				switch graph.Type {
				case "singlestat":
					components = append(components, singlestat(graph, data)...)
				case "gauge":
					components = append(components, gauge(graph, data)...)
				case "sparkline":
					components = append(components, sparkline(graph, data)...)
				case "plot":
					components = append(components, plot(graph, data)...)
				}
			}

			cols = append(cols, ui.NewCol(graph.Width, components...))
		}

		rows = append(rows, ui.NewRow(row.Height, cols...))
	}

	g.Set(rows...)
}

func singlestat(graph dashboard.Graph, data []datasource.Data) []interface{} {
	single := widgets.NewParagraph()
	single.Title = graph.Title

	var prefix string
	if graph.Options.Prefix != "" {
		prefix = graph.Options.Prefix + " "
	}

	var postfix string
	if graph.Options.Postfix != "" {
		postfix = " " + graph.Options.Postfix
	}

	if len(graph.Options.Stats) == 0 {
		graph.Options.Stats = []string{"current"}
	}

	var value string

	if len(data) == 0 {
		value = "N/A"
	} else {
		if graph.Options.Stats[0] == "name" {
			value = getLabelValue(graph.Options.Label, data[0].Labels)
		} else {
			floatValue := getStatValue(graph.Options.Stats[0], data[0].Points)
			if len(graph.Options.Thresholds) > 0 && len(graph.Options.Thresholds)+1 == len(graph.Options.Colors) {
				single.TextStyle.Fg = getColor(graph.Options.Colors[len(graph.Options.Colors)-1])
				for index, threshold := range graph.Options.Thresholds {
					if floatValue < threshold {
						single.TextStyle.Fg = getColor(graph.Options.Colors[index])
						break
					}
				}
			}

			value = strconv.FormatFloat(floatValue, 'f', graph.Options.Decimals, 64)
		}
	}

	single.Text = prefix + value + graph.Options.Unit + postfix
	single.WrapText = true

	components := make([]interface{}, 1)
	components[0] = single
	return components
}

func gauge(graph dashboard.Graph, data []datasource.Data) []interface{} {
	var components []interface{}

	if len(graph.Options.Stats) == 0 {
		graph.Options.Stats = []string{"current"}
	}

	for _, d := range data {
		g := widgets.NewGauge()

		title := graph.Title
		if graph.Options.Label != "" {
			title = title + " - " + getLabelValue(graph.Options.Label, d.Labels)
		}

		g.Title = title

		var value float64
		if len(d.Points) == 0 || math.IsNaN(d.Points[0]) {
			value = 0
		} else {
			value = getStatValue(graph.Options.Stats[0], d.Points)
			if len(graph.Options.Thresholds) > 0 && len(graph.Options.Thresholds)+1 == len(graph.Options.Colors) {
				g.BarColor = getColor(graph.Options.Colors[len(graph.Options.Colors)-1])
				for index, threshold := range graph.Options.Thresholds {
					if value < threshold {
						g.BarColor = getColor(graph.Options.Colors[index])
						break
					}
				}
			}
		}

		g.Percent = int(value)
		components = append(components, ui.NewRow(1.0/float64(len(data)), g))
	}

	return components
}

func sparkline(graph dashboard.Graph, data []datasource.Data) []interface{} {
	var sls []*widgets.Sparkline

	for _, d := range data {
		sl := widgets.NewSparkline()
		sl.Data = d.Points
		sl.Title = getLabelValue(graph.Options.Label, d.Labels)

		if len(d.Points) > 0 {
			sl.Title = sl.Title + ": " + strconv.FormatFloat(d.Points[len(d.Points)-1], 'f', graph.Options.Decimals, 64) + graph.Options.Unit

			var stats []string
			for _, stat := range graph.Options.Stats {
				stats = append(stats, stat+": "+strconv.FormatFloat(getStatValue(stat, d.Points), 'f', graph.Options.Decimals, 64))
			}

			if len(stats) > 0 {
				sl.Title = sl.Title + " (" + strings.Join(stats, ", ") + ")"
			}

			if len(graph.Options.Thresholds) > 0 && len(graph.Options.Thresholds)+1 == len(graph.Options.Colors) {
				sl.LineColor = getColor(graph.Options.Colors[len(graph.Options.Colors)-1])
				for index, threshold := range graph.Options.Thresholds {
					if d.Points[len(d.Points)-1] < threshold {
						sl.LineColor = getColor(graph.Options.Colors[index])
						break
					}
				}
			}
		}

		sls = append(sls, sl)
	}

	slg := widgets.NewSparklineGroup(sls...)
	slg.Title = graph.Title

	components := make([]interface{}, 1)
	components[0] = slg
	return components
}

func plot(graph dashboard.Graph, data []datasource.Data) []interface{} {
	plot := cPlot.NewPlot()
	plot.Title = graph.Title

	plotData := make([][]float64, len(data))
	var plotLabels []string
	for index, d := range data {
		plotData[index] = d.Points

		var stats []string
		for _, stat := range graph.Options.Stats {
			stats = append(stats, stat+": "+strconv.FormatFloat(getStatValue(stat, d.Points), 'f', graph.Options.Decimals, 64))
		}

		if len(stats) > 0 {
			plotLabels = append(plotLabels, getLabelValue(graph.Options.Label, d.Labels)+": "+strconv.FormatFloat(d.Points[len(d.Points)-1], 'f', graph.Options.Decimals, 64)+graph.Options.Unit+" ("+strings.Join(stats, ", ")+")")
		} else {
			plotLabels = append(plotLabels, getLabelValue(graph.Options.Label, d.Labels)+": "+strconv.FormatFloat(d.Points[len(d.Points)-1], 'f', graph.Options.Decimals, 64)+graph.Options.Unit)
		}

	}

	plot.Data = plotData
	plot.DataLabels = plotLabels

	components := make([]interface{}, 1)
	components[0] = plot
	return components
}

func getColor(color string) ui.Color {
	switch color {
	case "blue":
		return ui.ColorBlue
	case "cyan":
		return ui.ColorCyan
	case "green":
		return ui.ColorGreen
	case "magenta":
		return ui.ColorMagenta
	case "red":
		return ui.ColorRed
	case "white":
		return ui.ColorWhite
	case "yellow":
		return ui.ColorYellow
	default:
		return ui.ColorWhite
	}
}

func getLabelValue(label string, labels map[string]string) string {
	value, ok := labels[label]
	if !ok {
		var values []string
		for key, value := range labels {
			values = append(values, key+"="+value)
			return strings.Join(values, ", ")
		}
	}

	return value
}

func getStatValue(stat string, data []float64) float64 {
	switch stat {
	case "current":
		return data[len(data)-1]
	case "first":
		return data[0]
	case "min":
		min := data[0]
		for _, value := range data {
			if value < min {
				min = value
			}
		}
		return min
	case "max":
		max := data[0]
		for _, value := range data {
			if value > max {
				max = value
			}
		}
		return max
	case "avg":
		var total float64
		for _, value := range data {
			total = total + value
		}
		return total / float64(len(data))
	case "total":
		var total float64
		for _, value := range data {
			total = total + value
		}
		return total
	case "diff":
		return data[0] - data[len(data)]
	case "range":
		min := data[0]
		max := data[0]
		for _, value := range data {
			if value > max {
				max = value
			}
			if value < min {
				min = value
			}
		}
		return max - min
	default:
		return data[len(data)-1]
	}
}
