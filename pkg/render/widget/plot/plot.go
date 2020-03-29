// plot implements the plot chart from termui. This file is copied from the plot.go widgets from termui.
// See: https://github.com/gizak/termui/blob/master/v3/widgets/plot.go
//
// The following changes were made:
//   - Render legend in the top left corner of the chart
//   - Change the resolution of the data before rendering. This means we only render ever x point of the data if it
//     does does not fit into the draw area.
//
// Notes:
//   - Change the labels of the x axis to render the times of the selected time range.
//   - Missing data of a data series is ignored.
package plot

import (
	"fmt"
	"image"

	fLog "github.com/ricoberger/dash/pkg/log"

	. "github.com/gizak/termui/v3"
)

type Plot struct {
	Block

	Data       [][]float64
	DataLabels []string
	MaxVal     float64

	LineColors []Color
	AxesColor  Color
	ShowAxes   bool

	HorizontalScale int
}

const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 4
	xAxisLabelsGap    = 2
	yAxisLabelsGap    = 1
)

func NewPlot() *Plot {
	return &Plot{
		Block:           *NewBlock(),
		LineColors:      Theme.Plot.Lines,
		AxesColor:       Theme.Plot.Axes,
		Data:            [][]float64{},
		HorizontalScale: 1,
		ShowAxes:        true,
	}
}

func (self *Plot) renderBraille(buf *Buffer, drawArea image.Rectangle, maxVal float64) {
	canvas := NewCanvas()
	canvas.Rectangle = drawArea

	// Change the resolution of the data, so all data fits into the draw area.
	for i, line := range self.Data {
		resolution := float64(len(line)) / float64(drawArea.Max.X-drawArea.Min.X)
		fLog.Debugf("Resolution is %f, to fit the draw area from %d to %d", resolution, drawArea.Min.X, drawArea.Max.X)

		if int(resolution) > 1 {
			var data []float64
			for j := 0; j < len(line); j = j + int(resolution) {
				data = append(data, line[j])
			}

			fLog.Debugf("Change data resolution from %d to %d", len(line), len(data))
			self.Data[i] = data
		}
	}

	// Implementation of the original line chart from termui
	for i, line := range self.Data {
		previousHeight := int((line[1] / maxVal) * float64(drawArea.Dy()-1))
		for j, val := range line[1:] {
			height := int((val / maxVal) * float64(drawArea.Dy()-1))
			canvas.SetLine(
				image.Pt(
					(drawArea.Min.X+(j*self.HorizontalScale))*2,
					(drawArea.Max.Y-previousHeight-1)*4,
				),
				image.Pt(
					(drawArea.Min.X+((j+1)*self.HorizontalScale))*2,
					(drawArea.Max.Y-height-1)*4,
				),
				SelectColor(self.LineColors, i),
			)
			previousHeight = height
		}
	}

	canvas.Draw(buf)
}

func (self *Plot) plotLegend(buf *Buffer) {
	for index, label := range self.DataLabels {
		buf.SetString(
			label,
			NewStyle(self.LineColors[index]),
			image.Pt(self.Inner.Min.X+yAxisLabelsWidth+2, self.Inner.Min.Y+index),
		)
	}
}

func (self *Plot) plotAxes(buf *Buffer, maxVal float64) {
	// draw origin cell
	buf.SetCell(
		NewCell(BOTTOM_LEFT, NewStyle(ColorWhite)),
		image.Pt(self.Inner.Min.X+yAxisLabelsWidth, self.Inner.Max.Y-xAxisLabelsHeight-1),
	)
	// draw x axis line
	for i := yAxisLabelsWidth + 1; i < self.Inner.Dx(); i++ {
		buf.SetCell(
			NewCell(HORIZONTAL_DASH, NewStyle(ColorWhite)),
			image.Pt(i+self.Inner.Min.X, self.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}
	// draw y axis line
	for i := 0; i < self.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			NewCell(VERTICAL_DASH, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Min.X+yAxisLabelsWidth, i+self.Inner.Min.Y),
		)
	}
	// draw x axis labels
	// draw 0
	buf.SetString(
		"0",
		NewStyle(ColorWhite),
		image.Pt(self.Inner.Min.X+yAxisLabelsWidth, self.Inner.Max.Y-1),
	)
	// draw rest
	for x := self.Inner.Min.X + yAxisLabelsWidth + (xAxisLabelsGap)*self.HorizontalScale + 1; x < self.Inner.Max.X-1; {
		label := fmt.Sprintf(
			"%d",
			(x-(self.Inner.Min.X+yAxisLabelsWidth)-1)/(self.HorizontalScale)+1,
		)
		buf.SetString(
			label,
			NewStyle(ColorWhite),
			image.Pt(x, self.Inner.Max.Y-1),
		)
		x += (len(label) + xAxisLabelsGap) * self.HorizontalScale
	}

	// draw y axis labels
	verticalScale := maxVal / float64(self.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 0; i*(yAxisLabelsGap+1) < self.Inner.Dy()-1; i++ {
		buf.SetString(
			fmt.Sprintf("%.2f", float64(i)*verticalScale*(yAxisLabelsGap+1)),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Min.X, self.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2),
		)
	}
}

func (self *Plot) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	maxVal := self.MaxVal
	if maxVal == 0 {
		maxVal, _ = GetMaxFloat64From2dSlice(self.Data)
	}

	if self.ShowAxes {
		self.plotAxes(buf, maxVal)
	}

	drawArea := self.Inner
	if self.ShowAxes {
		drawArea = image.Rect(
			self.Inner.Min.X+yAxisLabelsWidth+1, self.Inner.Min.Y,
			self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
		)
	}

	self.renderBraille(buf, drawArea, maxVal)
	self.plotLegend(buf)
}
