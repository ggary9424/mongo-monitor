package termui

import (
	"context"
	"time"

	metricHelper "mongo-monitor/metric_helper"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/linechart"
)

// layoutButtons are buttons that change the layout.
type layoutButtons struct {
	lcB *button.Button
}

// widgets holds the widgets used by this demo.
type widgets struct {
	opcountersLC *linechart.LineChart
	buttons      *layoutButtons
}

// periodic executes the provided closure periodically every interval.
// Exits when the context expires.
func periodic(ctx context.Context, interval time.Duration, fn func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := fn(); err != nil {
				panic(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

var metricsSlice metricHelper.MetricsSlice

func UpdateMetricsSlice(ms metricHelper.MetricsSlice) {
	metricsSlice = ms
}

func extractOpcounters() (
	[]float64,
	[]float64,
	[]float64,
	[]float64,
	[]float64,
	[]float64,
	map[int]string,
) {
	var sortedMS metricHelper.MetricsSlice
	if metricsSlice == nil {
		return []float64{}, []float64{}, []float64{}, []float64{}, []float64{}, []float64{}, map[int]string{}
	} else {
		sortedMS = make(metricHelper.MetricsSlice, len(metricsSlice))
		copy(sortedMS, metricsSlice)
	}

	insertCountSlice := make([]float64, 50)
	queryCountSlice := make([]float64, 50)
	updateCountSlice := make([]float64, 50)
	deleteCountSlice := make([]float64, 50)
	getmoreCountSlice := make([]float64, 50)
	commandCountSlice := make([]float64, 50)
	index := 0
	for i := 50 - len(sortedMS); i < 50; i++ {
		insertCountSlice[i] = sortedMS[index].InsertCountPerSecond
		queryCountSlice[i] = sortedMS[index].QueryCountPerSecond
		updateCountSlice[i] = sortedMS[index].UpdateCountPerSecond
		deleteCountSlice[i] = sortedMS[index].DeleteCountPerSecond
		getmoreCountSlice[i] = sortedMS[index].GetmoreCountPerSecond
		commandCountSlice[i] = sortedMS[index].CommandCountPerSecond
		index++
	}
	XLabelMap := map[int]string{}
	return insertCountSlice, queryCountSlice, updateCountSlice, deleteCountSlice, getmoreCountSlice, commandCountSlice, XLabelMap
}

// newOpcounters returns a line chart that displays a heartbeat-like progression.
func newOpcounters(ctx context.Context) (*linechart.LineChart, error) {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
	)
	if err != nil {
		return nil, err
	}
	go periodic(ctx, redrawInterval/3, func() error {
		i, q, u, d, g, c, XLabelMap := extractOpcounters()
		err := lc.Series("insert", i,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(87))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		err = lc.Series("query", q,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(76))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		err = lc.Series("update", u,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(65))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		err = lc.Series("delete", d,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(54))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		err = lc.Series("getmore", g,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(43))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		err = lc.Series("command", c,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(32))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		return nil
	})
	return lc, nil
}

// newWidgets creates all widgets used by this demo.
func newWidgets(ctx context.Context, c *container.Container) (*widgets, error) {
	opcountersLC, err := newOpcounters(ctx)
	if err != nil {
		return nil, err
	}

	return &widgets{
		opcountersLC: opcountersLC,
	}, nil
}

// gridLayout prepares container options that represent the desired screen layout.
// This function demonstrates the use of the grid builder.
// gridLayout() and contLayout() demonstrate the two available layout APIs and
// both produce equivalent layouts for layoutType layoutAll.
func gridLayout(w *widgets) ([]container.Option, error) {
	builder := grid.New()
	builder.Add(
		grid.ColWidthPerc(99, grid.RowHeightPerc(99, grid.Widget(w.opcountersLC,
			container.Border(linestyle.Light),
			container.BorderTitle("opcounters"),
			container.BorderTitleAlignCenter(),
		))),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

const rootID = "root"

// redrawInterval is how often termdash redraws the screen.
const redrawInterval = 1 * time.Second

func Render(parentCtx context.Context) {
	t, err := termbox.New(termbox.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		panic(err)
	}
	defer t.Close()

	c, err := container.New(t, container.ID(rootID))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	w, err := newWidgets(ctx, c)
	if err != nil {
		panic(err)
	}

	gridOpts, err := gridLayout(w) // equivalent to contLayout(w)
	if err != nil {
		panic(err)
	}

	if err := c.Update(rootID, gridOpts...); err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC || k.Key.String() == "q" {
			cancel()
		}
	}
	if err := termdash.Run(
		ctx,
		t,
		c,
		termdash.KeyboardSubscriber(quitter),
		termdash.RedrawInterval(redrawInterval),
	); err != nil {
		panic(err)
	}
}
