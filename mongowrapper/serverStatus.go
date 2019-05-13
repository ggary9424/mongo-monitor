package mongowrapper

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// ServerStatus keeps the data returned by the serverStatus() method.
type ServerStatusStats struct {
	Version        string    `bson:"version"`
	Uptime         float64   `bson:"uptime"`
	UptimeEstimate float64   `bson:"uptimeEstimate"`
	LocalTime      time.Time `bson:"localTime"`

	// Asserts *AssertsStats `bson:"asserts"`
	Connections *ConnectionsStats `bson:"connections"`

	// Dur *DurStats `bson:"dur"`

	// BackgroundFlushing *FlushStats `bson:"backgroundFlushing"`

	// GlobalLock *GlobalLockStats `bson:"globalLock"`

	// IndexCounter *IndexCounterStats `bson:"indexCounters"`

	// Locks LockStatsMap `bson:"locks,omitempty"`

	Network *NetworkStats `bson:"network"`
	// OpLatencies    *OpLatenciesStat     `bson:"opLatencies"`
	Opcounters *OpcountersStats `bson:"opcounters"`
	// OpcountersRepl *OpcountersReplStats `bson:"opcountersRepl"`
	Metrics *MetricsStats `bson:"metrics"`

	// StorageEngine *StorageEngineStats `bson:"storageEngine"`
	// InMemory      *WiredTigerStats    `bson:"inMemory"`
	// RocksDb       *RocksDbStats       `bson:"rocksdb"`
	WiredTiger *WiredTigerStats `bson:"wiredTiger"`
}

// GetServerStatus returns the server status info.
func GetServerStatus(client *mongo.Client) *ServerStatusStats {
	serverStatus := &ServerStatusStats{}
	result := client.Database("admin").RunCommand(
		context.Background(),
		bsonx.Doc{
			{"serverStatus", bsonx.Int32(1)},
			{"recordStats", bsonx.Int32(0)},
			{"opLatencies", bsonx.Document(bsonx.MDoc{"histograms": bsonx.Boolean(true)})},
		},
	)
	result.Decode(serverStatus)
	return serverStatus
}
