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
	InsertCountPerSecondRecord     *CountPerSecondRecord
	QueryCountPerSecondRecord      *CountPerSecondRecord
	UpdateCountPerSecondRecord     *CountPerSecondRecord
	DeleteCountPerSecondRecord     *CountPerSecondRecord
	GetmoreCountPerSecondRecord    *CountPerSecondRecord
	CommandCountPerSecondRecord    *CountPerSecondRecord
	NetworkInBytesPerSecondRecord  *BytesPerSecondRecord
	NetworkOutBytesPerSecondRecord *BytesPerSecondRecord
	CheckpointCountPerSecondRecord *CountPerSecondRecord
}

var previousStatus *mongowrapper.ServerStatusStats

func ExtractMetrics(status *mongowrapper.ServerStatusStats) *Metrics {
	var metrics Metrics
	if previousStatus == nil {
		previousStatus = status
		return nil
	}
	metrics = Metrics{
		InsertCountPerSecondRecord:     getCountPerSecondByAction(ActionInsert, status),
		QueryCountPerSecondRecord:      getCountPerSecondByAction(ActionQuery, status),
		UpdateCountPerSecondRecord:     getCountPerSecondByAction(ActionUpdate, status),
		DeleteCountPerSecondRecord:     getCountPerSecondByAction(ActionDelete, status),
		GetmoreCountPerSecondRecord:    getCountPerSecondByAction(ActionGetmore, status),
		CommandCountPerSecondRecord:    getCountPerSecondByAction(ActionCommand, status),
		NetworkInBytesPerSecondRecord:  getBytesPerSecondByAction(DataNetworkIn, status),
		NetworkOutBytesPerSecondRecord: getBytesPerSecondByAction(DataNetworkOut, status),
		CheckpointCountPerSecondRecord: getCountPerSecondByAction(ActionCheckpoint, status),
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
