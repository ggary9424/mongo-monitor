package storage

import (
	metrichelper "mongo-monitor/metric_helper"
)

type Storage interface {
	FetchLastMetrics() (metrichelper.Metrics, error)
	FetchLastFewMetricsSlice(count int) (metrichelper.MetricsSlice, error)
	RecordMetrics(metrichelper.Metrics) error
}

type Driver int

const (
	Memory Driver = iota
)

func CreateStorage(driver Driver) Storage {
	switch driver {
	case Memory:
		return createMemoryStorage()
	default:
		return createMemoryStorage()
	}
}
