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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// "go.mongodb.org/mongo-driver/mongo"

type Con struct {
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
	//NewConnection(string) (*Con, error) // NewConnection-> constructor -> returns instance of Con, error
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

// NewConnection - constructor for our communications with Mongo
func NewConnection(dsn string, lg *slog.Logger) (*Con, error) {
	// create instance of Connection and log interface pointer
	ma := new(Con)
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

// Create -> create value with specified key on ourcollection
func (mc *Con) Create(collection, key string, value interface{}) (primitive.ObjectID, error) {
	mc.log.Debug("Creating new object in mongodb...Create")
	// first actualize collection
	mc.ActualCollection = mc.CLH.Database(mc.Database).Collection(collection)
	// now, we use insert1, to future we will make some logic above this package (probably in director)
	// don't want to use insert many, operation should be atomic! for multi insert we will make other one in future
	// but first, create and fill document metadata
<<<<<<< HEAD
	mdoc := newMongoDoc(collection, key, value)
=======
	mdoc := newMongoDoc(key, value)
	result, err := mc.ActualCollection.InsertOne(mc.CTX, mdoc)
	if err != nil {
		errstr := fmt.Sprintf("Error on inserting document: %v", err)
		mc.log.Error(errstr)
		return primitive.NilObjectID, err
	}
>>>>>>> 9129547 (2025-01-16 001)

	result, err := mc.ActualCollection.InsertOne(mc.CTX, mdoc)
	if err != nil {
		errstr := fmt.Sprintf("Error on inserting document: %v", err)
		return primitive.NilObjectID, errors.New(errstr)
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// Read -> read value specified by key
func (mc *Con) Read(collection, key string) (interface{}, error) {
	debugstr := fmt.Sprintf("Read from collection %s by key %s", collection, key)
	mc.log.Debug(debugstr)
	// first actualize collection
	mc.ActualCollection = mc.CLH.Database(mc.Database).Collection(collection)
	var result bson.M
	err := mc.ActualCollection.FindOne(mc.CTX, key).Decode(&result)
	if err != nil {
		errstr := fmt.Sprintf("Error on reading collection, key: %v", err)
		mc.log.Error(errstr)
		return nil, errors.New(errstr)
	}

	return result, nil
}

// Update -> check against old value, in case of difference update value
func (mc *Con) Update(collection, key string, value interface{}) (bool, error) {
	debugstr := fmt.Sprintf("Update data in collection %s identified by key %s", collection, key)
	mc.log.Debug(debugstr)
	// first actualize collection
	mc.ActualCollection = mc.CLH.Database(mc.Database).Collection(collection)
	// udpate document (as for https://www.mongodb.com/docs/drivers/go/current/quick-reference/ style)
	//find old values
	var firstResult interface{}
	err := mc.ActualCollection.FindOne(mc.CTX, key).Decode(firstResult)
	if err == mongo.ErrNoDocuments {
		errstr := fmt.Sprintf("No document with such key exists: %v", err)
		mc.log.Error(errstr)
		return false, errors.New(errstr)
	}
	if err != nil {
		errstr := fmt.Sprintf("Error on finding document to update: %v", err)
		mc.log.Error(errstr)
		return false, errors.New(errstr)
	}

	//var UpdateResult interface{}
	_, err = mc.ActualCollection.UpdateOne(mc.CTX, firstResult.(bson.D), bson.D{{key, value}})
	if err != nil {
		errstr := fmt.Sprintf("Error on updating document %v: %v", firstResult.(string), err)
		mc.log.Error(errstr)
		return false, errors.New(errstr)
	}

	return true, nil
}

func (mc *Con) Delete(collection, key string) error {
	debugstr := fmt.Sprintf("Delete from collection %s key %s", collection, key)
	mc.log.Debug(debugstr)
	mc.ActualCollection = mc.CLH.Database(mc.Database).Collection(collection)

	return nil
}

func (mc *Con) Close() error {
	return mc.CLH.Disconnect(mc.CTX)
}
