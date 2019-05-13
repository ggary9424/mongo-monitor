package cmd

import (
	"context"
	metrichelper "mongo-monitor/metric_helper"
	"mongo-monitor/mongowrapper"
	"mongo-monitor/storage"
	"os"
	"os/signal"
	"syscall"
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

func init() {
	rootCmd.AddCommand(mongostatCmd)
}

func mongostat() {
	var sigint chan os.Signal
	running := true
	go func() {
		sigint = make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from system
		signal.Notify(sigint, syscall.SIGTERM)

		switch <-sigint {
		case os.Interrupt:
			logrus.Warn("Receive single os.Interrupt")
		}

		running = false
	}()

	ctx := context.Background()
	defer ctx.Done()
	client, err := mongowrapper.CreateClient(viper.GetString("mongo.uri"))
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	s := storage.CreateStorage(storage.Memory)
	go recordMetricsPeriodically(ctx, client, s, 1*time.Second)
	go logMetricsPeriodically(ctx, s, 1*time.Second)
	for running {
		time.Sleep(1 * time.Second)
	}
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
			status := mongowrapper.GetServerStatus(client)
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
