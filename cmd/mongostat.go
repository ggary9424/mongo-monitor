package cmd

import (
	"context"
	"fmt"
	metrichelper "mongo-monitor/metric_helper"
	"mongo-monitor/mongowrapper"
	"mongo-monitor/storage"
	"mongo-monitor/termui"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongostatCmd will run mongostat function
var mongostatCmd = &cobra.Command{
	Use:   "mongostat",
	Short: "Just like mongostat command",
	Long:  "Just like mongostat command",
	Run: func(cmd *cobra.Command, args []string) {
		mongostat()
	},
}

var usingUI = false

func init() {
	pf := mongostatCmd.PersistentFlags()
	pf.BoolVar(&usingUI, "ui", false, "if you want to use UI or not")
	rootCmd.AddCommand(mongostatCmd)
}

func mongostat() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)

		signal.Notify(sigs, os.Interrupt)

		go func() {
			_, ok := <-sigs
			if ok {
				fmt.Println("Receive single os.Interrupt")
			} else {
				fmt.Println("Close the sigs channel")
			}
			done <- true
		}()

	Loop:
		for {
			select {
			case <-done:
				close(sigs)
				break Loop
			case <-ctx.Done():
				close(sigs)
				break Loop
			default:
			}
		}
	_:
		cancel()
		println("cancel")
	}()

	client, err := mongowrapper.CreateClient(ctx, viper.GetString("mongo.uri"))
	if err != nil {
		logrus.Error(err)
		panic(err)
	}

	s := storage.CreateStorage(storage.Memory)

	wg.Add(1)
	go func() {
		defer wg.Done()
		recordMetricsPeriodically(ctx, client, s, 1*time.Millisecond)
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if usingUI {
			go func() {
				termui.Render(ctx)
				cancel()
			}()
			updateTermuiDataPeriodically(ctx, s, 1*time.Millisecond)
		} else {
			logMetricsPeriodically(ctx, s, 1*time.Millisecond)
		}
		cancel()
	}()

	wg.Wait()
}

func recordMetricsPeriodically(
	ctx context.Context,
	client *mongo.Client,
	s storage.Storage,
	interval time.Duration,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			status := mongowrapper.GetServerStatus(ctx, client)
			metrics := metrichelper.ExtractMetrics(status)
			if metrics != nil {
				s.RecordMetrics(*metrics)
			}
			time.Sleep(interval)
		}
	}
}

func logMetricsPeriodically(
	ctx context.Context,
	s storage.Storage,
	interval time.Duration,
) error {
	logrus.Info("insert query update delete getmore command network_in network_out checkpoint")
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			metrics, err := s.FetchLastMetrics()
			if _, ok := err.(*storage.DataNotFound); ok {
				continue
			}
			logrus.Infof(
				"    *%d    *%d     *%d     *%d       %d       %d        %d       %d          %d\n",
				int64(metrics.InsertCountPerSecond),
				int64(metrics.QueryCountPerSecond),
				int64(metrics.UpdateCountPerSecond),
				int64(metrics.DeleteCountPerSecond),
				int64(metrics.GetmoreCountPerSecond),
				int64(metrics.CommandCountPerSecond),
				int64(metrics.NetworkInBytesPerSecond),
				int64(metrics.NetworkOutBytesPerSecond),
				// int64(status.Connections.Current),
				int64(metrics.CheckpointCountPerSecond),
			)
			time.Sleep(interval)
		}
	}
}

func updateTermuiDataPeriodically(
	ctx context.Context,
	s storage.Storage,
	interval time.Duration,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			ms, err := s.FetchLastFewMetricsSlice(50)
			if _, ok := err.(*storage.DataNotFound); ok {
				continue
			}
			termui.UpdateMetricsSlice(ms)
			time.Sleep(interval)
		}
	}
}
