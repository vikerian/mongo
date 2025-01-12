package mongo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/bson/primitive" -> driver using primitive.ObjectId

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// "go.mongodb.org/mongo-driver/mongo"

type MongoCon struct {
	log              *slog.Logger // our log interface
	URL              string
	CTX              context.Context
	Cancel           context.CancelFunc
	Options          *options.ClientOptions
	CLH              *mongo.Client
	Database         string
	ActualCollection *mongo.Collection
}

type MongoDB interface {
	//NewMongoConnection(string) (*MongoCon, error) // NewMongoConnection-> constructor -> returns instance of MongoCon, error
	Create(string, string, interface{}) (primitive.ObjectID, error) // Create(collection_name, key, value) -> objectid, error
	Read(string, string) (interface{}, error)                       // Read(collection_name, key, value) -> value, error
	Update(string, string, interface{}) (bool, error)               // Update(collection_name, key, value) -> ok, error
	Delete(string, string) (bool, error)                            // Delete(collection_name,key) -> ok,error
	Close() error                                                   // Cloase connection -> error
}

// ConstructDSN - create DSN for connection, mongo is quite picky on this
func MongoDBCreateDSN(username, password, host, port, database string) string {
	var mongoDSN string
	if username == "" || password == "" {
		mongoDSN = fmt.Sprintf("mongodb://%s:%s/%s", host, port, database)
	} else {
		//if username != "admin" {
		mongoDSN = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", username, password, host, port, database)
		//} else {
		//	mongoDSN = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", username, password, host, port, database)
	}
	return mongoDSN
}

// NewMongoConnection - constructor for our communications with Mongo
func NewMongoConnection(dsn string, lg *slog.Logger) (*MongoCon, error) {
	// create instance of MongoConnection and log interface pointer
	ma := new(MongoCon)
	if lg == nil {
		lg = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	ma.log.Info("Setting up mongo database connection...")
	// setup mongo conneciton
	ma.log = lg
	ma.URL = dsn
	ma.CTX, ma.Cancel = context.WithTimeout(context.Background(), 5*time.Second)
	// connect to db
	ma.Options = options.Client().ApplyURI(dsn)
	ma.log.Info("Connecting to mongo database...")
	clh, err := mongo.Connect(ma.CTX, ma.Options)
	ma.CLH = clh
	if err != nil {
		errstr := fmt.Sprintf("Error on setup connection to Mongo: %v", err)
		ma.log.Error(errstr)
		return nil, errors.New(errstr)
	}
	// actually now we touch db for first time
	err = ma.CLH.Ping(ma.CTX, nil)
	if err != nil {
		errstr := fmt.Sprintf("Error on checkup connection to Mongo: %v", err)
		ma.log.Error(errstr)
		return nil, errors.New(errstr)
	}
	ma.log.Info("Connection to mongodb successfull...")
	return ma, nil
}

// CreateVAL -> create value with specified key on ourcollection
func (mc *MongoCon) Create(collection, key string, value interface{}) (primitive.ObjectID, error) {
	mc.log.Debug("Creating new object in mongodb...Create")
	// first actualize collection
	mc.ActualCollection = mc.CLH.Database(mc.Database).Collection(collection)
	// now, we use insert1, to future we will make some logic above this package (probably in director)
	// don't want to use insert many, operation should be atomic! for multi insert we will make other one in future

	return primitive.NilObjectID, nil
}

// Read -> read value specified by key
func (mc *MongoCon) Read(collection, key string) (interface{}, error) {

	return nil, nil
}

// Update -> check against old value, in case of difference update value
func (mc *MongoCon) Update(collection, key string, value interface{}) (bool, error) {

	return true, nil
}

func (mc *MongoCon) Delete(collection, key string) error {

	return nil
}

func (mc *MongoCon) Close() error {
	return mc.CLH.Disconnect(mc.CTX)
}
