package widget

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/ricoberger/dash/pkg/dashboard"
	"github.com/ricoberger/dash/pkg/datasource"
	fLog "github.com/ricoberger/dash/pkg/log"
	"github.com/ricoberger/dash/pkg/render/utils"

	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/gauge"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/sparkline"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/olekukonko/tablewriter"
)

func GridLayout(storage *utils.Storage) []container.Option {
	var rows []grid.Element

	for _, row := range storage.Dashboard().Rows {
		var cols []grid.Element

		for _, graph := range row.Graphs {
			var component grid.Element

			if graph.Type == "table" {
				data, err := graph.GetTableData(storage.Datasource(), storage.VariableValues)
				if err != nil {
					component = renderError(graph, fmt.Sprintf("Could not load data: %s", err.Error()))
				} else {
					fLog.Debugf("TableData: %v", data)
					component, err = tablePanel(graph, data)
					if err != nil {
						component = renderError(graph, fmt.Sprintf("Could not render singlestat %s: %s", graph.Title, err.Error()))
					}
				}
			} else {
				data, err := graph.GetData(storage.Datasource(), storage.VariableValues, storage.Interval.Start, storage.Interval.End)
				if err != nil {
					component = renderError(graph, fmt.Sprintf("Could not load data: %s", err.Error()))
				} else {
					fLog.Debugf("render %d for %s", len(data.Series), graph.Title)

					switch graph.Type {
					case "singlestat":
						component, err = singlestatPanel(graph, data)
						if err != nil {
							component = renderError(graph, fmt.Sprintf("Could not render singlestat %s: %s", graph.Title, err.Error()))
						}
					case "gauge":
						component, err = gaugePanel(graph, data)
						if err != nil {
							component = renderError(graph, fmt.Sprintf("Could not render gauge %s: %s", graph.Title, err.Error()))
						}
					case "donut":
						component, err = donutPanel(graph, data)
						if err != nil {
							component = renderError(graph, fmt.Sprintf("Could not render donut %s: %s", graph.Title, err.Error()))
						}
					case "sparkline":
						component, err = sparklinePanel(graph, data)
						if err != nil {
							component = renderError(graph, fmt.Sprintf("Could not render sparkline %s: %s", graph.Title, err.Error()))
						}
					case "linechart":
						component, err = linechartPanel(graph, data)
						if err != nil {
							component = renderError(graph, fmt.Sprintf("Could not load render linechart %s: %s", graph.Title, err.Error()))
						}
					}
				}
			}

			cols = append(cols, grid.ColWidthPerc(graph.Width, component))
		}

		rows = append(rows, grid.RowHeightPerc(row.Height, cols...))
	}

	builder := grid.New()
	builder.Add(rows...)
	gridOpts, _ := builder.Build()
	return gridOpts
}

func renderError(graph dashboard.Graph, err string) grid.Element {
	log.Printf(err)
	txt, _ := text.New()
	txt.Write(err)

	return grid.Widget(
		txt,
		container.Border(linestyle.Light),
		container.BorderTitle(graph.Title),
		container.AlignHorizontal(align.HorizontalCenter),
		container.AlignVertical(align.VerticalMiddle),
	)
}

func singlestatPanel(graph dashboard.Graph, data *datasource.Data) (grid.Element, error) {
	single, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}

	if len(graph.Options.Stats) == 0 {
		graph.Options.Stats = []string{"current"}
	}

	var value string
	var color cell.Color

	if len(data.Series) == 0 {
		value = "NaN"
	} else {
		if graph.Options.Stats[0] == "name" {
			value = data.Series[0].Label
		} else {
			floatValue := getStatValue(graph.Options.Stats[0], data.Series[0].Points)
			if len(graph.Options.Thresholds) > 0 && len(graph.Options.Thresholds)+1 == len(graph.Options.Colors) {
				color = getColor(graph.Options.Colors[len(graph.Options.Colors)-1])
				for index, threshold := range graph.Options.Thresholds {
					if floatValue < threshold {
						color = getColor(graph.Options.Colors[index])
						break
					}
				}
			}

			value = strconv.FormatFloat(floatValue, 'f', graph.Options.Decimals, 64)
		}
	}

	if mapping, ok := graph.Options.Mappings[value]; ok {
		value = mapping
	}

	var chunks []*segmentdisplay.TextChunk
	chunk := segmentdisplay.NewChunk(value+" "+graph.Options.Unit, segmentdisplay.WriteCellOpts(cell.FgColor(color)))
	chunks = append(chunks, chunk)

	err = single.Write(chunks)
	if err != nil {
		return nil, err
	}

	return grid.Widget(single, container.Border(linestyle.Light), container.BorderTitle(graph.Title), container.AlignHorizontal(align.HorizontalCenter), container.AlignVertical(align.VerticalMiddle)), nil
}

func gaugePanel(graph dashboard.Graph, data *datasource.Data) (grid.Element, error) {
	if len(graph.Options.Stats) == 0 {
		graph.Options.Stats = []string{"current"}
	}

	var value float64
	var color cell.Color

	if len(data.Series) == 0 || math.IsNaN(data.Series[0].Points[0]) {
		value = 0
	} else {
		value = getStatValue(graph.Options.Stats[0], data.Series[0].Points)
		if len(graph.Options.Thresholds) > 0 && len(graph.Options.Thresholds)+1 == len(graph.Options.Colors) {
			color = getColor(graph.Options.Colors[len(graph.Options.Colors)-1])
			for index, threshold := range graph.Options.Thresholds {
				if value < threshold {
					color = getColor(graph.Options.Colors[index])
					break
				}
			}
		}
	}

	g, err := gauge.New(gauge.Color(color))
	if err != nil {
		return nil, err
	}

	err = g.Percent(int(value))
	if err != nil {
		return nil, err
	}

	return grid.Widget(g, container.Border(linestyle.Light), container.BorderTitle(graph.Title), container.AlignHorizontal(align.HorizontalCenter), container.AlignVertical(align.VerticalMiddle)), nil
}

func donutPanel(graph dashboard.Graph, data *datasource.Data) (grid.Element, error) {
	if len(graph.Options.Stats) == 0 {
		graph.Options.Stats = []string{"current"}
	}

	var value float64
	var color cell.Color

	if len(data.Series) == 0 || math.IsNaN(data.Series[0].Points[0]) {
		value = 0
	} else {
		value = getStatValue(graph.Options.Stats[0], data.Series[0].Points)
		if len(graph.Options.Thresholds) > 0 && len(graph.Options.Thresholds)+1 == len(graph.Options.Colors) {
			color = getColor(graph.Options.Colors[len(graph.Options.Colors)-1])
			for index, threshold := range graph.Options.Thresholds {
				if value < threshold {
					color = getColor(graph.Options.Colors[index])
					break
				}
			}
		}
	}

	d, err := donut.New(donut.CellOpts(cell.FgColor(color)))
	if err != nil {
		return nil, err
	}

	err = d.Percent(int(value))
	if err != nil {
		return nil, err
	}

	return grid.Widget(d, container.Border(linestyle.Light), container.BorderTitle(graph.Title), container.AlignHorizontal(align.HorizontalCenter), container.AlignVertical(align.VerticalMiddle)), nil
}

func sparklinePanel(graph dashboard.Graph, data *datasource.Data) (grid.Element, error) {
	var values []int
	var color cell.Color
	var label string

	if len(data.Series) > 0 {
		for _, value := range data.Series[0].Points {
			values = append(values, int(value))
		}

		if len(graph.Options.Thresholds) > 0 && len(graph.Options.Thresholds)+1 == len(graph.Options.Colors) {
			color = getColor(graph.Options.Colors[len(graph.Options.Colors)-1])
			for index, threshold := range graph.Options.Thresholds {
				if values[len(values)-1] < int(threshold) {
					color = getColor(graph.Options.Colors[index])
					break
				}
			}
		}

		label = fmt.Sprintf("%s: %s %s", data.Series[0].Label, strconv.FormatFloat(data.Series[0].Points[len(data.Series[0].Points)-1], 'f', graph.Options.Decimals, 64), graph.Options.Unit)
	}

	s, err := sparkline.New(sparkline.Label(label, cell.FgColor(color)), sparkline.Color(color))
	if err != nil {
		return nil, err
	}

	err = s.Add(values)
	if err != nil {
		return nil, err
	}

	return grid.Widget(s, container.Border(linestyle.Light), container.BorderTitle(graph.Title), container.AlignHorizontal(align.HorizontalCenter), container.AlignVertical(align.VerticalMiddle)), nil
}

func linechartPanel(graph dashboard.Graph, data *datasource.Data) (grid.Element, error) {
	lc, err := linechart.New()
	if err != nil {
		return nil, err
	}

	legend, err := text.New(text.WrapAtRunes())
	if err != nil {
		return nil, err
	}

	for index, series := range data.Series {
		var stats []string
		for _, stat := range graph.Options.Stats {
			stats = append(stats, fmt.Sprintf("%s: %s", stat, strconv.FormatFloat(getStatValue(stat, series.Points), 'f', graph.Options.Decimals, 64)))
		}

		var statsLegend string
		if len(stats) > 0 {
			statsLegend = fmt.Sprintf("%s %s (%s)", strconv.FormatFloat(getStatValue("current", series.Points), 'f', graph.Options.Decimals, 64), graph.Options.Unit, strings.Join(stats, ", "))
		} else {
			statsLegend = fmt.Sprintf("%s %s", strconv.FormatFloat(getStatValue("current", series.Points), 'f', graph.Options.Decimals, 64), graph.Options.Unit)
		}

		color := randomColor(index)
		if graph.Options.Legend == "bottom" {
			err = legend.Write(fmt.Sprintf("%s: %s   ", series.Label, statsLegend), text.WriteCellOpts(cell.FgColor(color)))
			if err != nil {
				return nil, err
			}
		} else if graph.Options.Legend == "right" {
			err = legend.Write(fmt.Sprintf("%s: %s\n", series.Label, statsLegend), text.WriteCellOpts(cell.FgColor(color)))
			if err != nil {
				return nil, err
			}
		}

		if index == 0 {
			err = lc.Series(series.Label, series.Points, linechart.SeriesCellOpts(cell.FgColor(color)), linechart.SeriesXLabels(data.Timestamps))
			if err != nil {
				return nil, err
			}
		} else {
			err = lc.Series(series.Label, series.Points, linechart.SeriesCellOpts(cell.FgColor(color)))
			if err != nil {
				return nil, err
			}
		}
	}

	// Render linechart and legend
	// See: https://github.com/slok/grafterm/blob/master/internal/view/render/termdash/graph.go
	graphElement := grid.Widget(lc)

	var elements []grid.Element
	switch graph.Options.Legend {
	case "bottom":
		legendElement := grid.RowHeightPercWithOpts(99, []container.Option{container.PaddingTopPercent(50)}, grid.Widget(legend))
		elements = []grid.Element{grid.RowHeightPerc(90, graphElement), grid.RowHeightPerc(4, legendElement)}
	case "right":
		legendElement := grid.ColWidthPercWithOpts(99, []container.Option{container.PaddingLeftPercent(10)}, grid.Widget(legend))
		elements = []grid.Element{grid.ColWidthPerc(80, graphElement), grid.ColWidthPerc(19, legendElement)}
	default:
		elements = []grid.Element{grid.ColWidthPerc(99, graphElement)}
	}

	opts := []container.Option{container.Border(linestyle.Light), container.BorderTitle(graph.Title), container.AlignHorizontal(align.HorizontalCenter), container.AlignVertical(align.VerticalMiddle)}
	element := grid.RowHeightPercWithOpts(99, opts, elements...)

	return element, nil
}

func tablePanel(graph dashboard.Graph, data *datasource.TableData) (grid.Element, error) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	var headers []string
	var names []string

	for _, column := range graph.Options.Columns {
		headers = append(headers, column.Header)
		names = append(names, column.Name)
	}

	for _, value := range *data {
		var columns []string

		for _, name := range names {
			columns = append(columns, formateInterface(value[name], graph.Options.Decimals))
		}

		table.Append(columns)
	}

	table.SetHeader(headers)
	table.Render()

	txt, err := text.New()
	if err != nil {
		return nil, err
	}

	err = txt.Write(tableString.String())
	if err != nil {
		return nil, err
	}

	return grid.Widget(txt, container.Border(linestyle.Light), container.BorderTitle(graph.Title), container.AlignHorizontal(align.HorizontalCenter), container.AlignVertical(align.VerticalMiddle)), nil
}

func formateInterface(value interface{}, decimals int) string {
	switch i := value.(type) {
	case float64:
		return strconv.FormatFloat(i, 'f', decimals, 64)
	case string:
		return i
	default:
		fLog.Debugf("Could not formate value: %v, type: %v", i, reflect.TypeOf(i))
		return fmt.Sprintf("%v", i)
	}
}

func getColor(color string) cell.Color {
	switch color {
	case "blue":
		return cell.ColorBlue
	case "cyan":
		return cell.ColorCyan
	case "green":
		return cell.ColorGreen
	case "magenta":
		return cell.ColorMagenta
	case "red":
		return cell.ColorRed
	case "white":
		return cell.ColorWhite
	case "yellow":
		return cell.ColorYellow
	default:
		return cell.ColorWhite
	}
}

func randomColor(index int) cell.Color {
	var colors = []cell.Color{cell.ColorBlue, cell.ColorCyan, cell.ColorGreen, cell.ColorMagenta, cell.ColorRed, cell.ColorWhite, cell.ColorYellow}

	if index < len(colors) {
		return colors[index]
	}

	return cell.ColorNumber(rand.Intn(255-0) + 0)
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
