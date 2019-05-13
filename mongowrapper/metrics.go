package mongowrapper

// DocumentStats are the stats associated to a document.
type DocumentStats struct {
	Deleted  float64 `bson:"deleted"`
	Inserted float64 `bson:"inserted"`
	Returned float64 `bson:"returned"`
	Updated  float64 `bson:"updated"`
}

// BenchmarkStats is bechmark info about an operation.
type BenchmarkStats struct {
	Num         float64 `bson:"num"`
	TotalMillis float64 `bson:"totalMillis"`
}

// GetLastErrorStats are the last error stats.
type GetLastErrorStats struct {
	Wtimeouts float64         `bson:"wtimeouts"`
	Wtime     *BenchmarkStats `bson:"wtime"`
}

// OperationStats are the stats for some kind of operations.
type OperationStats struct {
	Fastmod      float64 `bson:"fastmod"`
	Idhack       float64 `bson:"idhack"`
	ScanAndOrder float64 `bson:"scanAndOrder"`
}

// QueryExecutorStats are the stats associated with a query execution.
type QueryExecutorStats struct {
	Scanned        float64 `bson:"scanned"`
	ScannedObjects float64 `bson:"scannedObjects"`
}

// RecordStats are stats associated with a record.
type RecordStats struct {
	Moves float64 `bson:"moves"`
}

// ApplyStats are the stats associated with the apply operation.
type ApplyStats struct {
	Batches *BenchmarkStats `bson:"batches"`
	Ops     float64         `bson:"ops"`
}

// BufferStats are the stats associated with the buffer
type BufferStats struct {
	Count        float64 `bson:"count"`
	MaxSizeBytes float64 `bson:"maxSizeBytes"`
	SizeBytes    float64 `bson:"sizeBytes"`
}

// ReplExecutorStats are the stats associated with replication execution
type ReplExecutorStats struct {
	Counters         map[string]float64 `bson:"counters"`
	Queues           map[string]float64 `bson:"queues"`
	EventWaiters     float64            `bson:"eventWaiters"`
	UnsignaledEvents float64            `bson:"unsignaledEvents"`
}

// MetricsNetworkStats are the network stats.
type MetricsNetworkStats struct {
	Bytes          float64         `bson:"bytes"`
	Ops            float64         `bson:"ops"`
	GetMores       *BenchmarkStats `bson:"getmores"`
	ReadersCreated float64         `bson:"readersCreated"`
}

// ReplStats are the stats associated with the replication process.
type ReplStats struct {
	Apply        *ApplyStats          `bson:"apply"`
	Buffer       *BufferStats         `bson:"buffer"`
	Executor     *ReplExecutorStats   `bson:"executor,omitempty"`
	Network      *MetricsNetworkStats `bson:"network"`
	PreloadStats *PreloadStats        `bson:"preload"`
}

// PreloadStats are the stats associated with preload operation.
type PreloadStats struct {
	Docs    *BenchmarkStats `bson:"docs"`
	Indexes *BenchmarkStats `bson:"indexes"`
}

// StorageStats are the stats associated with the storage.
type StorageStats struct {
	BucketExhausted float64 `bson:"freelist.search.bucketExhausted"`
	Requests        float64 `bson:"freelist.search.requests"`
	Scanned         float64 `bson:"freelist.search.scanned"`
}

// CursorStatsOpen are the stats for open cursors
type CursorStatsOpen struct {
	NoTimeout float64 `bson:"noTimeout"`
	Pinned    float64 `bson:"pinned"`
	Total     float64 `bson:"total"`
}

// CursorStats are the stats for cursors
type CursorStats struct {
	TimedOut float64          `bson:"timedOut"`
	Open     *CursorStatsOpen `bson:"open"`
}

// TTLStats are the stats for ttl indexes
type TTLStats struct {
	DeletedDocuments float64 `bson:"deletedDocuments"`
	Passes           float64 `bson:"passes"`
}

// MetricsStats are all stats associated with metrics of the system
type MetricsStats struct {
	Document      *DocumentStats      `bson:"document"`
	GetLastError  *GetLastErrorStats  `bson:"getLastError"`
	Operation     *OperationStats     `bson:"operation"`
	QueryExecutor *QueryExecutorStats `bson:"queryExecutor"`
	Record        *RecordStats        `bson:"record"`
	Repl          *ReplStats          `bson:"repl"`
	Storage       *StorageStats       `bson:"storage"`
	Cursor        *CursorStats        `bson:"cursor"`
	TTL           *TTLStats           `bson:"ttl"`
}
