package mongowrapper

// ConnectionsStats are connections metrics
type ConnectionsStats struct {
	Current      float64 `bson:"current"`
	Available    float64 `bson:"available"`
	TotalCreated float64 `bson:"totalCreated"`
}
