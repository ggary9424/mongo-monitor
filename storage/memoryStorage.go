package storage

import (
	metrichelper "mongo-monitor/metric_helper"
	"sync"
)

type MemoryStorage struct{}

type recordsWithMutex struct {
	records []metrichelper.Metrics
	mutex   sync.Mutex
}

var recordsWM = recordsWithMutex{
	records: []metrichelper.Metrics{},
}

type DataNotFound struct{}

func (e *DataNotFound) Error() string {
	return "Data not found"
}

func (storage *MemoryStorage) FetchLastMetrics() (metrichelper.Metrics, error) {
	recordsWM.mutex.Lock()
	records := recordsWM.records
	recordsWM.mutex.Unlock()
	if len(records) < 1 {
		return metrichelper.Metrics{}, &DataNotFound{}
	}
	return records[len(records)-1], nil
}

func (storage *MemoryStorage) FetchLastFewMetricsSlice(count int) (metrichelper.MetricsSlice, error) {
	recordsWM.mutex.Lock()
	records := recordsWM.records
	recordsWM.mutex.Unlock()
	if len(records) < 1 {
		return metrichelper.MetricsSlice{}, &DataNotFound{}
	}
	start := (map[bool]int{true: len(records) - count, false: 0})[len(records)-count > 0]
	return records[start:], nil
}

func (storage *MemoryStorage) RecordMetrics(metrics metrichelper.Metrics) error {
	recordsWM.mutex.Lock()
	recordsWM.records = append(recordsWM.records, metrics)
	recordsWM.mutex.Unlock()
	return nil
}

func createMemoryStorage() Storage {
	return &MemoryStorage{}
}
