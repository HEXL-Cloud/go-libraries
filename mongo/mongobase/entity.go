package mongobase

// This struct is only used to demonstrate the usage of the MongoBaseRepository,
// and is not intended to be used directly in the application.
//
// The consumer should define their own entity struct with the required fields
type Entity struct {
	ID string `bson:"_id,omitempty"`
}
