package mongo

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UnixTimestamp is a custom type that helps to Marshal a Datetime value
// from MongoDB to Unix epoch timestamp in the JSON response
type UnixTimestamp time.Time

func (t *UnixTimestamp) Time() time.Time {
	return time.Time(*t)
}

// MarshalJSON helps to parse Datetime from database to unix epoch timestamp
func (t UnixTimestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(t).Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

// UnmarshalJSON cast time.Time input to UnixTimestamp type
func (t *UnixTimestamp) UnmarshalJSON(b []byte) error {
	var ti = time.Time(*t)
	res := ti.UnmarshalJSON(b)
	*t = UnixTimestamp(ti)

	return res
}

// GetBSON helps to store timestamp as Datetime in MongoDB
func (t UnixTimestamp) GetBSON() (interface{}, error) {
	return time.Time(t), nil
}

// SetBSON creates a UnixTimestamp from a MongoDB datetime string.
func (t *UnixTimestamp) SetBSON(raw bson.RawValue) error {
	// create datetime type object to reuse its UnmarshalJSON method
	datetime := primitive.NewDateTimeFromTime(t.Time())

	err := raw.Unmarshal(&datetime)
	// set back unmarshalled value
	*t = UnixTimestamp(datetime.Time())

	return err
}
