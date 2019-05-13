package metric_helper

import (
	"mongo-monitor/mongowrapper"
	"time"

	"github.com/sirupsen/logrus"
)

type CommandCountPerSecondRecord struct {
	StartTime time.Time
	EndTime   time.Time
	Count     float64
}

var records []CommandCountPerSecondRecord

func recordCommandCountPerSecond(status *mongowrapper.ServerStatusStats) CommandCountPerSecondRecord {
	previousUnixTime := previousStatus.LocalTime
	currentUnixTime := status.LocalTime
	previousInsertCount := previousStatus.Opcounters.Command
	currentInsertCount := status.Opcounters.Command
	insertCountPerSecond := (currentInsertCount - previousInsertCount) / float64(currentUnixTime.Unix()-previousUnixTime.Unix())
	logrus.Info(previousUnixTime, currentUnixTime, previousInsertCount, currentInsertCount, insertCountPerSecond)

	record := CommandCountPerSecondRecord{
		StartTime: previousUnixTime,
		EndTime:   currentUnixTime,
		Count:     insertCountPerSecond,
	}
	records = append(records, record)

	return record
}

func GetCommandCountPerSecondRecords() []CommandCountPerSecondRecord {
	return records
}
