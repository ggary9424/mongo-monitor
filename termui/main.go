package termui

import (
	"context"
	"fmt"
	"time"

	metricHelper "mongo-monitor/metric_helper"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
)

// layoutButtons are buttons that change the layout.
type layoutButtons struct {
	lcB *button.Button
}

// widgets holds the widgets used by this demo.
type widgets struct {
	mongostatUIText *text.Text
	opcountersLC    *linechart.LineChart
	opcountersText  *text.Text
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

// newOpcountersLc returns a line chart that displays opcounters line chart.
func newOpcountersLc(ctx context.Context) (*linechart.LineChart, error) {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorNumber(161))),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorNumber(222))),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorNumber(222))),
		linechart.XAxisUnscaled(),
	)
	if err != nil {
		return nil, err
	}
	go periodic(ctx, redrawInterval/3, func() error {
		i, q, u, d, g, c, XLabelMap := extractOpcounters()
		err := lc.Series("insert", i,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(111))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		err = lc.Series("query", q,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(172))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		err = lc.Series("update", u,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(107))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		err = lc.Series("delete", d,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(161))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		err = lc.Series("getmore", g,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(245))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		err = lc.Series("command", c,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(135))),
			linechart.SeriesXLabels(XLabelMap),
		)
		if err != nil {
			return err
		}
		return nil
	})
	return lc, nil
}

// newMongostatChartText returns a text block that displays the basic infomation of mongostat chart.
func newMongostatUIText(ctx context.Context) (*text.Text, error) {
	t, err := text.New()
	if err != nil {
		return nil, err
	}
	if err := t.Write(
		fmt.Sprintf("Press Esc/Q/Ctrl-C to quit\n"),
		text.WriteCellOpts(cell.FgColor(cell.ColorNumber(111))),
	); err != nil {
		return nil, err
	}

	return t, nil
}

// newOpcountersText returns a text block that displays the infomation of opcounters line chart.
func newOpcountersText(ctx context.Context) (*text.Text, error) {
	t, err := text.New()
	if err != nil {
		return nil, err
	}

	if err := t.Write(fmt.Sprintf("INSERT  .......\n"), text.WriteCellOpts(cell.FgColor(cell.ColorNumber(111)))); err != nil {
		return nil, err
	}
	if err := t.Write(fmt.Sprintf("QUERY   .......\n"), text.WriteCellOpts(cell.FgColor(cell.ColorNumber(172)))); err != nil {
		return nil, err
	}
	if err := t.Write(fmt.Sprintf("UPDATE  .......\n"), text.WriteCellOpts(cell.FgColor(cell.ColorNumber(107)))); err != nil {
		return nil, err
	}
	if err := t.Write(fmt.Sprintf("DELETE  .......\n"), text.WriteCellOpts(cell.FgColor(cell.ColorNumber(161)))); err != nil {
		return nil, err
	}
	if err := t.Write(fmt.Sprintf("GETMORE .......\n"), text.WriteCellOpts(cell.FgColor(cell.ColorNumber(245)))); err != nil {
		return nil, err
	}
	if err := t.Write(fmt.Sprintf("COMMAND .......\n"), text.WriteCellOpts(cell.FgColor(cell.ColorNumber(135)))); err != nil {
		return nil, err
	}

	return t, nil
}

// newWidgets creates all widgets used by this demo.
func newWidgets(ctx context.Context, c *container.Container) (*widgets, error) {
	mongostatUIText, err := newMongostatUIText(ctx)
	if err != nil {
		return nil, err
	}

	opcountersLC, err := newOpcountersLc(ctx)
	if err != nil {
		return nil, err
	}

	opcountersText, err := newOpcountersText(ctx)
	if err != nil {
		return nil, err
	}

	return &widgets{
		mongostatUIText: mongostatUIText,
		opcountersLC:    opcountersLC,
		opcountersText:  opcountersText,
	}, nil
}

func getLayoutOpts(w *widgets) ([]container.Option, error) {
	return []container.Option{
		container.SplitHorizontal(
			container.Top(
				container.PlaceWidget(w.mongostatUIText),
				container.Border(linestyle.Light),
				container.BorderTitle("Welcome to Mongostat UI"),
				container.BorderTitleAlignCenter(),
			),
			container.Bottom(
				container.SplitVertical(
					container.Left(
						container.PlaceWidget(w.opcountersText),
						container.Border(linestyle.Light),
						container.BorderTitle("Lines"),
						container.BorderTitleAlignCenter(),
					),
					container.Right(
						container.PlaceWidget(w.opcountersLC),
						container.Border(linestyle.Light),
						container.BorderTitle("Opcounters Line Chart"),
						container.BorderTitleAlignCenter(),
					),
					container.SplitPercent(15),
				),
			),
			container.SplitPercent(12),
		),
	}, nil
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

	layoutOpts, err := getLayoutOpts(w)
	if err != nil {
		panic(err)
	}

	if err := c.Update(rootID, layoutOpts...); err != nil {
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
