package metric_helper

import (
	"mongo-monitor/mongowrapper"
	"time"
)

type ActionType int

const (
	ActionInsert ActionType = iota
	ActionQuery
	ActionUpdate
	ActionDelete
	ActionGetmore
	ActionCommand
	ActionCheckpoint
)

type DataType int

const (
	DataNetworkIn DataType = iota
	DataNetworkOut
)

type CountPerSecondRecord struct {
	ActionType ActionType
	StartTime  time.Time
	EndTime    time.Time
	Count      float64
}

type BytesPerSecondRecord struct {
	DataType  DataType
	StartTime time.Time
	EndTime   time.Time
	Bytes     float64
}

type Metrics struct {
	InsertCountPerSecond     float64
	QueryCountPerSecond      float64
	UpdateCountPerSecond     float64
	DeleteCountPerSecond     float64
	GetmoreCountPerSecond    float64
	CommandCountPerSecond    float64
	NetworkInBytesPerSecond  float64
	NetworkOutBytesPerSecond float64
	CheckpointCountPerSecond float64
	StartTime                time.Time
	EndTime                  time.Time
}

var previousStatus *mongowrapper.ServerStatusStats

func ExtractMetrics(status *mongowrapper.ServerStatusStats) *Metrics {
	var metrics Metrics
	if previousStatus == nil {
		previousStatus = status
		return nil
	}
	metrics = Metrics{
		InsertCountPerSecond:     getCountPerSecondByAction(ActionInsert, status).Count,
		QueryCountPerSecond:      getCountPerSecondByAction(ActionQuery, status).Count,
		UpdateCountPerSecond:     getCountPerSecondByAction(ActionUpdate, status).Count,
		DeleteCountPerSecond:     getCountPerSecondByAction(ActionDelete, status).Count,
		GetmoreCountPerSecond:    getCountPerSecondByAction(ActionGetmore, status).Count,
		CommandCountPerSecond:    getCountPerSecondByAction(ActionCommand, status).Count,
		NetworkInBytesPerSecond:  getBytesPerSecondByAction(DataNetworkIn, status).Bytes,
		NetworkOutBytesPerSecond: getBytesPerSecondByAction(DataNetworkOut, status).Bytes,
		CheckpointCountPerSecond: getCountPerSecondByAction(ActionCheckpoint, status).Count,
		StartTime:                previousStatus.LocalTime,
		EndTime:                  status.LocalTime,
	}
	previousStatus = status
	return &metrics
}

func getCountPerSecondByAction(
	actionType ActionType,
	status *mongowrapper.ServerStatusStats,
) *CountPerSecondRecord {
	previousTime := previousStatus.LocalTime
	currentTime := status.LocalTime
	var previousCount, currentCount float64

	switch actionType {
	case ActionInsert:
		previousCount = previousStatus.Opcounters.Insert
		currentCount = status.Opcounters.Insert
	case ActionQuery:
		previousCount = previousStatus.Opcounters.Query
		currentCount = status.Opcounters.Query
	case ActionUpdate:
		previousCount = previousStatus.Opcounters.Update
		currentCount = status.Opcounters.Update
	case ActionDelete:
		previousCount = previousStatus.Opcounters.Delete
		currentCount = status.Opcounters.Delete
	case ActionGetmore:
		previousCount = previousStatus.Opcounters.GetMore
		currentCount = status.Opcounters.GetMore
	case ActionCommand:
		previousCount = previousStatus.Opcounters.Command
		currentCount = status.Opcounters.Command
	case ActionCheckpoint:
		previousCount = previousStatus.WiredTiger.Transaction.Checkpoints
		currentCount = status.WiredTiger.Transaction.Checkpoints
	}
	countPerSecond := (currentCount - previousCount) / float64(currentTime.Unix()-previousTime.Unix())

	return &CountPerSecondRecord{
		ActionType: actionType,
		StartTime:  previousTime,
		EndTime:    currentTime,
		Count:      countPerSecond,
	}
}

func getBytesPerSecondByAction(
	dataType DataType,
	status *mongowrapper.ServerStatusStats,
) *BytesPerSecondRecord {
	previousTime := previousStatus.LocalTime
	currentTime := status.LocalTime
	var previousBytes, currentBytes float64

	switch dataType {
	case DataNetworkIn:
		previousBytes = previousStatus.Network.BytesIn
		currentBytes = status.Network.BytesIn
	case DataNetworkOut:
		previousBytes = previousStatus.Network.BytesOut
		currentBytes = status.Network.BytesOut
	}
	bytesPerSecond := (currentBytes - previousBytes) / float64(currentTime.Unix()-previousTime.Unix())

	return &BytesPerSecondRecord{
		DataType:  dataType,
		StartTime: previousTime,
		EndTime:   currentTime,
		Bytes:     bytesPerSecond,
	}
}
