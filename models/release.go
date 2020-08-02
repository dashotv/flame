package models

import (
	"time"

	"github.com/Kamva/mgm/v3"
	"github.com/Kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Release struct {
	Document `bson:",inline"` // include mgm.DefaultModel
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type        string    `json:"type" bson:"type"`
	Source      string    `json:"source" bson:"source"`
	Raw         string    `json:"raw" bson:"raw"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Size        string    `json:"size" bson:"size"`
	View        string    `json:"view" bson:"view"`
	Download    string    `json:"download" bson:"download"`
	Infohash    string    `json:"infohash" bson:"infohash"`
	Name        string    `json:"name" bson:"name"`
	Season      int       `json:"season" bson:"season"`
	Episode     int       `json:"episode" bson:"episode"`
	Volume      int       `json:"volume" bson:"volume"`
	Checksum    string    `json:"checksum" bson:"checksum"`
	Group       string    `json:"group" bson:"group"`
	Author      string    `json:"author" bson:"author"`
	Verified    bool      `json:"verified" bson:"verified"`
	Widescreen  bool      `json:"widescreen" bson:"widescreen"`
	Uncensored  bool      `json:"uncensored" bson:"uncensored"`
	Bluray      bool      `json:"bluray" bson:"bluray"`
	Resolution  string    `json:"resolution" bson:"resolution"`
	Encoding    string    `json:"encoding" bson:"encoding"`
	Quality     string    `json:"quality" bson:"quality"`
	Published   time.Time `json:"published" bson:"published_at"`
}

func NewRelease() *Release {
	return &Release{}
}

type ReleaseStore struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mgm.Collection
}

func NewReleaseStore(URI, db, name string) (*ReleaseStore, error) {
	client, err := mgm.NewClient(CustomClientOptions(URI))
	if err != nil {
		return nil, err
	}

	database := client.Database(db)
	collection := mgm.NewCollection(database, name)

	store := &ReleaseStore{
		Client:     client,
		Database:   database,
		Collection: collection,
	}

	return store, nil
}

func (s *ReleaseStore) FindByID(id primitive.ObjectID) (*Release, error) {
	c := NewRelease()
	err := s.Collection.FindByID(id, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *ReleaseStore) Find(id string) (*Release, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.FindByID(oid)
}

func (s *ReleaseStore) Save(o *Release) error {
	// TODO: if id is nil create otherwise, call update
	return s.Collection.Create(o)
}

func (s *ReleaseStore) Update(o *Release) error {
	return s.Collection.Update(o)
}

func (s *ReleaseStore) Delete(o *Release) error {
	return s.Collection.Delete(o)
}

func (s *ReleaseStore) Query() *ReleaseQuery {
	values := make(bson.M)
	return &ReleaseQuery{
		store:  s,
		values: values,
		limit:  25,
		skip:   0,
		sort:   bson.D{},
	}
}

type ReleaseQuery struct {
	store  *ReleaseStore
	values bson.M
	limit  int64
	skip   int64
	sort   bson.D
}

func (q *ReleaseQuery) addSort(field string, value int) *ReleaseQuery {
	q.sort = append(q.sort, bson.E{Key: field, Value: value})
	return q
}

func (q *ReleaseQuery) Asc(field string) *ReleaseQuery {
	return q.addSort(field, 1)
}

func (q *ReleaseQuery) Desc(field string) *ReleaseQuery {
	return q.addSort(field, -1)
}

func (q *ReleaseQuery) Limit(limit int) *ReleaseQuery {
	q.limit = int64(limit)
	return q
}

func (q *ReleaseQuery) Skip(skip int) *ReleaseQuery {
	q.skip = int64(skip)
	return q
}

func (q *ReleaseQuery) options() *options.FindOptions {
	o := &options.FindOptions{}
	o.SetLimit(q.limit)
	o.SetSkip(q.skip)
	o.SetSort(q.sort)
	return o
}

func (q *ReleaseQuery) Run() ([]Release, error) {
	result := make([]Release, 0)
	err := q.store.Collection.SimpleFind(&result, q.values, q.options())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (q *ReleaseQuery) Where(key string, value interface{}) *ReleaseQuery {
	q.values[key] = bson.M{operator.Eq: value}
	return q
}

func (q *ReleaseQuery) In(key string, value interface{}) *ReleaseQuery {
	q.values[key] = bson.M{operator.In: value}
	return q
}

func (q *ReleaseQuery) NotIn(key string, value interface{}) *ReleaseQuery {
	q.values[key] = bson.M{operator.Nin: value}
	return q
}

func (q *ReleaseQuery) NotEqual(key string, value interface{}) *ReleaseQuery {
	q.values[key] = bson.M{operator.Ne: value}
	return q
}

func (q *ReleaseQuery) LessThan(key string, value interface{}) *ReleaseQuery {
	q.values[key] = bson.M{operator.Lt: value}
	return q
}

func (q *ReleaseQuery) LessThanEqual(key string, value interface{}) *ReleaseQuery {
	q.values[key] = bson.M{operator.Lte: value}
	return q
}

func (q *ReleaseQuery) GreaterThan(key string, value interface{}) *ReleaseQuery {
	q.values[key] = bson.M{operator.Gt: value}
	return q
}

func (q *ReleaseQuery) GreaterThanEqual(key string, value interface{}) *ReleaseQuery {
	q.values[key] = bson.M{operator.Gte: value}
	return q
}
