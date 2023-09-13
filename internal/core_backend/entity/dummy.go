package entity

type Dummy struct {
	ID      int    `bson:"id"`
	Message string `bson:"message"`
}

// CollectionName Collection name of Dummy
func (Dummy) CollectionName() string {
	return "dummies"
}
