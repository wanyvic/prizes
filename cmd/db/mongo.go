package db

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types/swarm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var (
	mongoDBReader  *mongo.Client
	ctxWithTimeout = context.Background()
	MongDBClient   = mongodb{uri: "mongodb://localhost:27017", database: "docker"}
)

type mongodb struct {
	uri      string `json:",omitempty"`
	database string `json:",omitempty"`
}

func (m *mongodb) InsertServiceOne(service swarm.Service) (bool, error) {
	if err := RefreshMongoDBConnection(mongoDBReader, &m.uri); err != nil {
		return false, err
	}
	collection := mongoDBReader.Database(m.database).Collection("service")
	_, err := collection.InsertOne(ctxWithTimeout, service)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *mongodb) InsertTaskOne(task swarm.Task) (bool, error) {
	if err := RefreshMongoDBConnection(mongoDBReader, &m.uri); err != nil {
		return false, err
	}
	collection := mongoDBReader.Database(m.database).Collection("task")
	_, err := collection.InsertOne(ctxWithTimeout, task)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *mongodb) UpdateServiceOne(service swarm.Service) (bool, error) {
	if err := RefreshMongoDBConnection(mongoDBReader, &m.uri); err != nil {
		return false, err
	}
	updateOption := options.UpdateOptions{}
	updateOption.SetUpsert(true)
	collection := mongoDBReader.Database(m.database).Collection("service")
	_, err := collection.UpdateOne(ctxWithTimeout, bson.D{{"id", service.ID}}, bson.D{{"$set", service}}, &updateOption)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *mongodb) UpdateTaskOne(task swarm.Task) (bool, error) {
	if err := RefreshMongoDBConnection(mongoDBReader, &m.uri); err != nil {
		return false, err
	}
	updateOption := options.UpdateOptions{}
	updateOption.SetUpsert(true)
	collection := mongoDBReader.Database(m.database).Collection("task")
	_, err := collection.UpdateOne(ctxWithTimeout, bson.D{{"id", task.ID}}, bson.D{{"$set", task}}, &updateOption)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *mongodb) FindServiceOne(serviceID string) (*swarm.Service, error) {
	var service swarm.Service
	if err := RefreshMongoDBConnection(mongoDBReader, &m.uri); err != nil {
		return nil, err
	}
	collection := mongoDBReader.Database(m.database).Collection("service")

	if err := collection.FindOne(ctxWithTimeout, bson.D{{"id", serviceID}}).Decode(&service); err != nil {
		return nil, err
	}
	return &service, nil
}
func (m *mongodb) FindTaskList(serviceID string) (*[]swarm.Task, error) {
	var tasklist []swarm.Task
	if err := RefreshMongoDBConnection(mongoDBReader, &m.uri); err != nil {
		return nil, err
	}
	collection := mongoDBReader.Database(m.database).Collection("task")

	cursor, err := collection.Find(ctxWithTimeout, bson.D{{"serviceid", serviceID}}, nil)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctxWithTimeout) {
		var task swarm.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasklist = append(tasklist, task)
	}
	return &tasklist, nil
}
func init() {
	// Handle Terminate Signal(Ctrl + C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = disconnectMongoDBClient(mongoDBReader)
	}()
}

// dialMongoDB connects a client to MongoDB server.
// Returns a MongoDB Client or any dialing error.
func dialMongoDB(uri *string) (*mongo.Client, error) {
	if strings.TrimSpace(*uri) == "" {
		return nil, errors.New("uri format error")
	}
	client, err := mongo.Connect(ctxWithTimeout, options.Client().ApplyURI(*uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}
	return client, nil
}

// disconnectMongoDBClient disconnects a client from MongoDB server.
// Returns if there is any disconnection error.
func disconnectMongoDBClient(client *mongo.Client) error {
	if client == nil {
		return errors.New("ErrNilMongoDBClient")
	}
	return client.Disconnect(ctxWithTimeout)
}

// refreshMongoDBConnection refreshes a client's connection with MongoDB server.
// Returns if there is any connection error.
func RefreshMongoDBConnection(client *mongo.Client, uri *string) error {
	if client == nil {
		newClient, err := dialMongoDB(uri)
		if err != nil {
			return err
		}
		mongoDBReader = newClient
		return nil
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		newClient, err := dialMongoDB(uri)
		if err != nil {
			mongoDBReader = nil
			return err
		}
		mongoDBReader = newClient
		return nil
	}
	return nil
}
