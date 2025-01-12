package mongo

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoDoc struct {
	Id         uuid.UUID          `json:"uuid"`
	MongoID    primitive.ObjectID `json:"mongo_id,omitempty"`
	Collection string             `json:"mongo_collection,omitempty"`
	Name       string             `json:"doc_name"`
	Version    string             `json:"doc_version,omitempty"`
	Authors    []Author           `json:"authors,omitempty"`
	Raw        []byte             `json:"raw,omitempty"`
	Document   interface{}        `json:"doc_structured,omitempty"`
	Data       []interface{}      `json:"data,omitempty"`
}

type Author struct {
	Id      uuid.UUID `json:"author_id,omitempty"`
	Name    string    `json:"author_name,omitempty"`
	Contact []string  `json:"author_contact_informations,omitempty"`
}
