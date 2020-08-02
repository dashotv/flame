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

type Download struct {
	Document `bson:",inline"` // include mgm.DefaultModel
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	MediumId   primitive.ObjectID `json:"medium_id" bson:"medium_id"`
	Auto       bool               `json:"auto" bson:"auto"`
	Multi      bool               `json:"multi" bson:"multi"`
	Force      bool               `json:"force" bson:"force"`
	Url        string             `json:"url" bson:"url"`
	ReleaseId  string             `json:"release_id" bson:"tdo_id"`
	Thash      string             `json:"thash" bson:"thash"`
	Timestamps struct {
		Found      time.Time `json:"found" bson:"found"`
		Loaded     time.Time `json:"loaded" bson:"loaded"`
		Downloaded time.Time `json:"downloaded" bson:"downloaded"`
		Completed  time.Time `json:"completed" bson:"completed"`
		Deleted    time.Time `json:"deleted" bson:"deleted"`
	} `json:"timestamps" bson:"timestamps"`
	Selected string `json:"selected" bson:"selected"`
	Status   string `json:"status" bson:"status"`
	Files    []struct {
		Id       primitive.ObjectID `json:"id" bson:"_id"`
		MediumId primitive.ObjectID `json:"medium_id" bson:"medium_id"`
		Num      int                `json:"num" bson:"num"`
	} `json:"download_files" bson:"download_files"`
}

func NewDownload() *Download {
	return &Download{}
}

type DownloadStore struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mgm.Collection
}

func NewDownloadStore(URI, db, name string) (*DownloadStore, error) {
	client, err := mgm.NewClient(CustomClientOptions(URI))
	if err != nil {
		return nil, err
	}

	database := client.Database(db)
	collection := mgm.NewCollection(database, name)

	store := &DownloadStore{
		Client:     client,
		Database:   database,
		Collection: collection,
	}

	return store, nil
}

func (s *DownloadStore) FindByID(id primitive.ObjectID) (*Download, error) {
	c := NewDownload()
	err := s.Collection.FindByID(id, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *DownloadStore) Find(id string) (*Download, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.FindByID(oid)
}

func (s *DownloadStore) Save(o *Download) error {
	// TODO: if id is nil create otherwise, call update
	return s.Collection.Create(o)
}

func (s *DownloadStore) Update(o *Download) error {
	return s.Collection.Update(o)
}

func (s *DownloadStore) Delete(o *Download) error {
	return s.Collection.Delete(o)
}

func (s *DownloadStore) Query() *DownloadQuery {
	values := make(bson.M)
	return &DownloadQuery{
		store:  s,
		values: values,
		limit:  25,
		skip:   0,
		sort:   bson.D{},
	}
}

type DownloadQuery struct {
	store  *DownloadStore
	values bson.M
	limit  int64
	skip   int64
	sort   bson.D
}

func (q *DownloadQuery) addSort(field string, value int) *DownloadQuery {
	q.sort = append(q.sort, bson.E{Key: field, Value: value})
	return q
}

func (q *DownloadQuery) Asc(field string) *DownloadQuery {
	return q.addSort(field, 1)
}

func (q *DownloadQuery) Desc(field string) *DownloadQuery {
	return q.addSort(field, -1)
}

func (q *DownloadQuery) Limit(limit int) *DownloadQuery {
	q.limit = int64(limit)
	return q
}

func (q *DownloadQuery) Skip(skip int) *DownloadQuery {
	q.skip = int64(skip)
	return q
}

func (q *DownloadQuery) options() *options.FindOptions {
	o := &options.FindOptions{}
	o.SetLimit(q.limit)
	o.SetSkip(q.skip)
	o.SetSort(q.sort)
	return o
}

func (q *DownloadQuery) Run() ([]Download, error) {
	result := make([]Download, 0)
	err := q.store.Collection.SimpleFind(&result, q.values, q.options())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (q *DownloadQuery) Where(key string, value interface{}) *DownloadQuery {
	q.values[key] = bson.M{operator.Eq: value}
	return q
}

func (q *DownloadQuery) In(key string, value interface{}) *DownloadQuery {
	q.values[key] = bson.M{operator.In: value}
	return q
}

func (q *DownloadQuery) NotIn(key string, value interface{}) *DownloadQuery {
	q.values[key] = bson.M{operator.Nin: value}
	return q
}

func (q *DownloadQuery) NotEqual(key string, value interface{}) *DownloadQuery {
	q.values[key] = bson.M{operator.Ne: value}
	return q
}

func (q *DownloadQuery) LessThan(key string, value interface{}) *DownloadQuery {
	q.values[key] = bson.M{operator.Lt: value}
	return q
}

func (q *DownloadQuery) LessThanEqual(key string, value interface{}) *DownloadQuery {
	q.values[key] = bson.M{operator.Lte: value}
	return q
}

func (q *DownloadQuery) GreaterThan(key string, value interface{}) *DownloadQuery {
	q.values[key] = bson.M{operator.Gt: value}
	return q
}

func (q *DownloadQuery) GreaterThanEqual(key string, value interface{}) *DownloadQuery {
	q.values[key] = bson.M{operator.Gte: value}
	return q
}
