package mongo

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoDoc struct {
	id          uuid.UUID          `bson:"uuid"`
	mongoID     primitive.ObjectID `bson:"mongo_id,omitempty"`
	collection  string             `bson:"mongo_collection,omitempty"`
	name        string             `bson:"mongo_name"`
	version     string             `bson:"doc_version,omitempty"`
	authors     []Author           `bson:"doc_authors,omitempty"`
	raw         []byte             `bson:"raw_bytes,omitempty"`
	document    interface{}        `bson:"doc,omitempty"`
	data        []interface{}      `bson:"data,omitempty"`
	ctime       time.Time          `bson:"mongo_doc_creation_time"`
	atimes      []time.Time        `bson:"mongo_doc_access_times,omitempty"`
	utimes      []time.Time        `bson:"mongo_doc_update_times,omitempty"`
	validity    time.Time          `bson:"mongo_doc_valid_till,omitempty"`
	deprecation time.Time          `bson:"mongo_doc_deprecation,omitempty"`
}

type Author struct {
	authorId     uuid.UUID `bson:"author_id,omitempty"`
	authorName   string    `bson:"author_name,omitempty"`
	authoContact []string  `bson:"author_contact_informations,omitempty"`
}

// NewMongoDoc -> create instance of mongo document with basic data filled in
func newMongoDoc(docname string, document interface{}) *mongoDoc {
	mdoc := new(mongoDoc)
	mdoc.id = uuid.New()
	mdoc.name = docname
	mdoc.document = document
	mdoc.ctime = time.Now()
	return mdoc
}

// addRaw -> add raw data ([]byte)
func (mdoc *mongoDoc) addRAW(rawdata []byte) {
	mdoc.raw = rawdata
}

// add version
func (mdoc *mongoDoc) addVersion(version string) {
	mdoc.version = version
}
