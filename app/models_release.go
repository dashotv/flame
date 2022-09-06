package app

import (
	"time"

	"github.com/dashotv/grimoire"
)

type Release struct {
	grimoire.Document `bson:",inline"` // includes default model settings
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
