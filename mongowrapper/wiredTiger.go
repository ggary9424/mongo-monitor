package mongowrapper

// WTTransactionStats stats
type WTTransactionStats struct {
	Begins               float64 `bson:"transaction begins"`
	Checkpoints          float64 `bson:"transaction checkpoints"`
	CheckpointsRunning   float64 `bson:"transaction checkpoint currently running"`
	CheckpointMaxMs      float64 `bson:"transaction checkpoint max time (msecs)"`
	CheckpointMinMs      float64 `bson:"transaction checkpoint min time (msecs)"`
	CheckpointLastMs     float64 `bson:"transaction checkpoint most recent time (msecs)"`
	CheckpointTotalMs    float64 `bson:"transaction checkpoint total time (msecs)"`
	Committed            float64 `bson:"transactions committed"`
	CacheOverflowFailure float64 `bson:"transaction failures due to cache overflow"`
	RolledBack           float64 `bson:"transactions rolled back"`
}

// WiredTiger stats
type WiredTigerStats struct {
	// BlockManager           *WTBlockManagerStats           `bson:"block-manager"`
	// Cache                  *WTCacheStats                  `bson:"cache"`
	// Log                    *WTLogStats                    `bson:"log"`
	// Session                *WTSessionStats                `bson:"session"`
	Transaction *WTTransactionStats `bson:"transaction"`
	// ConcurrentTransactions *WTConcurrentTransactionsStats `bson:"concurrentTransactions"`
}
