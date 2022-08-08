package models

import (
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Medium struct {
	Document `bson:",inline"` // include mgm.DefaultModel
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type         string           `json:"type" bson:"_type"`
	Kind         primitive.Symbol `json:"kind" bson:"kind"`
	Source       string           `json:"source" bson:"source"`
	SourceId     string           `json:"source_id" bson:"source_id"`
	Title        string           `json:"title" bson:"title"`
	Description  string           `json:"description" bson:"description"`
	Slug         string           `json:"slug" bson:"slug"`
	Text         []string         `json:"text" bson:"text"`
	Display      string           `json:"display" bson:"display"`
	Directory    string           `json:"directory" bson:"directory"`
	Search       string           `json:"search" bson:"search"`
	SearchParams struct {
		Type       string `json:"type" bson:"type"`
		Verified   bool   `json:"verified" bson:"verified"`
		Group      string `json:"group" bson:"group"`
		Author     string `json:"author" bson:"author"`
		Resolution int    `json:"resolution" bson:"resolution"`
		Source     string `json:"source" bson:"source"`
		Uncensored bool   `json:"uncensored" bson:"uncensored"`
		Bluray     bool   `json:"bluray" bson:"bluray"`
	} `json:"search_params" bson:"search_params"`
	Active      bool      `json:"active" bson:"active"`
	Downloaded  bool      `json:"downloaded" bson:"downloaded"`
	Completed   bool      `json:"completed" bson:"completed"`
	Skipped     bool      `json:"skipped" bson:"skipped"`
	Watched     bool      `json:"watched" bson:"watched"`
	Broken      bool      `json:"broken" bson:"broken"`
	ReleaseDate time.Time `json:"release_date" bson:"release_date"`
	Paths       []struct {
		Id        primitive.ObjectID `json:"ID" bson:"ID"`
		Type      primitive.Symbol   `json:"type" bson:"type"`
		Remote    string             `json:"remote" bson:"remote"`
		Local     string             `json:"local" bson:"local"`
		Extension string             `json:"extension" bson:"extension"`
		Size      int                `json:"size" bson:"size"`
		UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	} `json:"paths" bson:"paths"`
}

func NewMedium() *Medium {
	return &Medium{}
}

type MediumStore struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mgm.Collection
}

func NewMediumStore(URI, db, name string) (*MediumStore, error) {
	client, err := mgm.NewClient(CustomClientOptions(URI))
	if err != nil {
		return nil, err
	}

	database := client.Database(db)
	collection := mgm.NewCollection(database, name)

	store := &MediumStore{
		Client:     client,
		Database:   database,
		Collection: collection,
	}

	return store, nil
}

func (s *MediumStore) FindByID(id primitive.ObjectID) (*Medium, error) {
	c := NewMedium()
	err := s.Collection.FindByID(id, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *MediumStore) Find(id string) (*Medium, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.FindByID(oid)
}

func (s *MediumStore) Save(o *Medium) error {
	// TODO: if id is nil create otherwise, call update
	return s.Collection.Create(o)
}

func (s *MediumStore) Update(o *Medium) error {
	return s.Collection.Update(o)
}

func (s *MediumStore) Delete(o *Medium) error {
	return s.Collection.Delete(o)
}

func (s *MediumStore) Query() *MediumQuery {
	values := make(bson.M)
	return &MediumQuery{
		store:  s,
		values: values,
		limit:  25,
		skip:   0,
		sort:   bson.D{},
	}
}

type MediumQuery struct {
	store  *MediumStore
	values bson.M
	limit  int64
	skip   int64
	sort   bson.D
}

func (q *MediumQuery) addSort(field string, value int) *MediumQuery {
	q.sort = append(q.sort, bson.E{Key: field, Value: value})
	return q
}

func (q *MediumQuery) Asc(field string) *MediumQuery {
	return q.addSort(field, 1)
}

func (q *MediumQuery) Desc(field string) *MediumQuery {
	return q.addSort(field, -1)
}

func (q *MediumQuery) Limit(limit int) *MediumQuery {
	q.limit = int64(limit)
	return q
}

func (q *MediumQuery) Skip(skip int) *MediumQuery {
	q.skip = int64(skip)
	return q
}

func (q *MediumQuery) options() *options.FindOptions {
	o := &options.FindOptions{}
	o.SetLimit(q.limit)
	o.SetSkip(q.skip)
	o.SetSort(q.sort)
	return o
}

func (q *MediumQuery) Run() ([]Medium, error) {
	result := make([]Medium, 0)
	err := q.store.Collection.SimpleFind(&result, q.values, q.options())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (q *MediumQuery) Where(key string, value interface{}) *MediumQuery {
	q.values[key] = bson.M{operator.Eq: value}
	return q
}

func (q *MediumQuery) In(key string, value interface{}) *MediumQuery {
	q.values[key] = bson.M{operator.In: value}
	return q
}

func (q *MediumQuery) NotIn(key string, value interface{}) *MediumQuery {
	q.values[key] = bson.M{operator.Nin: value}
	return q
}

func (q *MediumQuery) NotEqual(key string, value interface{}) *MediumQuery {
	q.values[key] = bson.M{operator.Ne: value}
	return q
}

func (q *MediumQuery) LessThan(key string, value interface{}) *MediumQuery {
	q.values[key] = bson.M{operator.Lt: value}
	return q
}

func (q *MediumQuery) LessThanEqual(key string, value interface{}) *MediumQuery {
	q.values[key] = bson.M{operator.Lte: value}
	return q
}

func (q *MediumQuery) GreaterThan(key string, value interface{}) *MediumQuery {
	q.values[key] = bson.M{operator.Gt: value}
	return q
}

func (q *MediumQuery) GreaterThanEqual(key string, value interface{}) *MediumQuery {
	q.values[key] = bson.M{operator.Gte: value}
	return q
}
