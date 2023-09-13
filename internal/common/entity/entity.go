package cmentity

import (
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"backend-service/pkg/common/utils"
)

// Entity define base fields for all struct entity.
type Entity struct {
	ID        ID        `json:"id" bson:"_id"`
	CreatedAt CreatedAt `json:"-" bson:"created_at"`
	UpdatedAt CreatedAt `json:"-" bson:"updated_at"`
	Status    string    `json:"status" bson:"status"`
}

// ID is a custom type that helps to Marshal id value from Database.
// to string in the JSON response.
type ID string

// ToID returns ID type from string.
func ToID(id string) *ID {
	return utils.ToPtr(ID(id))
}

// String returns string of ID type.
func (id ID) String() string {
	return string(id)
}

// NewID returns new random ID.
func NewID() ID {
	return ID(primitive.NilObjectID.Hex())
}

// GetBSON helps to store ID as primitive.ObjectID in database.
//
// specific use for mongoDB.
func (id ID) GetBSON() (any, error) {
	if id.String() == "" {
		return primitive.NewObjectID(), nil
	}

	objID, err := primitive.ObjectIDFromHex(id.String())

	if err != nil {
		return nil, err
	}

	return objID, nil
}

func NewIDFromHexString(hexStr string) (*ID, error) {
	if hexStr == "" {
		return nil, nil
	}

	objectID, err := primitive.ObjectIDFromHex(hexStr)
	if err != nil {

		return nil, err
	}

	return utils.ToPtr(ID(objectID.Hex())), nil
}

func (id *ID) SetBSON(raw bson.RawValue) error {
	var idVal string

	if err := raw.Unmarshal(&idVal); err != nil {
		return err
	}

	*id = ID(idVal)

	return nil
}

// CreatedAt is a custom type that helps to Marshal a Datetime value from Database.
// to Unix epoch timestamp in the JSON response and auto upsert time when upsert document.
type CreatedAt time.Time

// MarshalJSON helps to parse Datetime from database to unix epoch timestamp
func (c CreatedAt) MarshalJSON() ([]byte, error) {
	ts := c.Time().Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

// UnmarshalJSON cast time.Time or unix epoch timestamp input to UnixTimestamp type.
func (c *CreatedAt) UnmarshalJSON(b []byte) error {
	unix, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		ti := time.Time(*c)
		res := ti.UnmarshalJSON(b)
		*c = CreatedAt(ti)

		return res
	}

	ti := time.Unix(unix, 0)
	*c = CreatedAt(ti)

	return nil
}

// Time returns time.Time value.
func (c *CreatedAt) Time() time.Time {
	unixTime := UnixTimestamp(*c)

	return unixTime.Time()
}

// GetBSON will run when inserting or updating createdAt field to DB, this helps to auto set time.
func (c CreatedAt) GetBSON() (any, error) {
	if c.Time().IsZero() {
		return time.Now(), nil
	}

	return c.Time(), nil
}

// SetBSON will run to decode BSON raw value.
func (c *CreatedAt) SetBSON(raw bson.RawValue) error {
	// create datetime type object to reuse its UnmarshalJSON method
	datetime := primitive.NewDateTimeFromTime(c.Time())

	err := raw.Unmarshal(&datetime)
	if err == nil {
		// set back unmarshalled value
		*c = CreatedAt(datetime.Time())

		return nil
	}

	return nil
}

// UnixTimestamp is a custom type that helps to Marshal a Datetime value from Database.
// to Unix epoch timestamp in the JSON response.
type UnixTimestamp time.Time

// Time returns time.Time value.
func (t *UnixTimestamp) Time() time.Time {
	return time.Time(*t)
}

// MarshalJSON helps to parse Datetime from database to unix epoch timestamp.
func (t UnixTimestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(t).Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

// UnmarshalJSON cast time.Time or unix epoch timestamp input to UnixTimestamp type.
func (t *UnixTimestamp) UnmarshalJSON(b []byte) error {
	unix, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		ti := time.Time(*t)
		res := ti.UnmarshalJSON(b)
		*t = UnixTimestamp(ti)

		return res
	}

	var ti = time.Unix(unix, 0)
	*t = UnixTimestamp(ti)

	return nil
}

// GetBSON helps to store timestamp as Datetime in database.
func (t UnixTimestamp) GetBSON() (interface{}, error) {
	return time.Time(t), nil
}

// SetBSON creates a UnixTimestamp from a MongoDB datetime string.
func (t *UnixTimestamp) SetBSON(raw bson.RawValue) error {
	// create datetime type object to reuse its UnmarshalJSON method
	datetime := primitive.NewDateTimeFromTime(t.Time())

	err := raw.Unmarshal(&datetime)
	if err == nil {
		// set back unmarshalled value
		*t = UnixTimestamp(datetime.Time())

		return nil
	}

	var unix int64

	err = raw.Unmarshal(&unix)
	if err != nil {
		return err
	}

	var ti = time.Unix(unix, 0)
	*t = UnixTimestamp(ti)

	return nil
}

type Media struct {
	URL          string `bson:"url" json:"url"`
	Type         string `bson:"type" json:"type"`
	ThumbnailURL string `bson:"thumbnail_url" json:"thumbnail_url"`
}
