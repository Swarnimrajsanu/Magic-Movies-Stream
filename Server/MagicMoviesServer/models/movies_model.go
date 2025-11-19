package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genere struct {
	GenereID   int    `bson:"genere_id" json:"genere_id" validate:"required"`
	GenereName string `bson:"genere_name" json:"genere_name" validate:"required,min=1,max=100"`
}
type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required"`
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"required"`
}

type Movie struct {
	ID     bson.ObjectID `bson:"_id" json:"id"`
	ImdbID string        `bson:"imdb_id" json:"imdb_id" validate:"required"`

	Title       string   `bson:"title" json:"title" validate:"required,min=1,max=500"`
	Posterpath  string   `bson:"poster_path" json:"poster_path" validate:"required,url"`
	YouTubeID   string   `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genere      []Genere `bson:"genere" json:"genere" validate:"required,dive"`
	AdminReview string   `bson:"admin_review" json:"admin_review" validate:"required,min=10,max=2000"`
	Ranking     Ranking  `bson:"ranking" json:"ranking" validate:"required"`
}
