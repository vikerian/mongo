package mongo

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoDoc struct {
	Id         uuid.UUID          `bson:"uuid"`
	MongoID    primitive.ObjectID `bson:"mongo_id,omitempty"`
	Collection string             `bson:"mongo_collection,omitempty"`
	Name       string             `bson:"doc_name"`
	Version    string             `bson:"doc_version,omitempty"`
	Authors    []Author           `bson:"authors,omitempty"`
	Raw        []byte             `bson:"raw,omitempty"`
	Document   interface{}        `bson:"doc_structured,omitempty"`
	Data       []interface{}      `bson:"data,omitempty"`
}

type Author struct {
	Id      uuid.UUID `bson:"author_id,omitempty"`
	Name    string    `bson:"author_name,omitempty"`
	Contact []string  `bson:"author_contact_informations,omitempty"`
}
