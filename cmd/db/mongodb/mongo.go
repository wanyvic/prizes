package mongodb

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types/swarm"
	"github.com/wanyvic/prizes/api/types/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var (
	MongoDBDefaultURI      = "mongodb://localhost:27017"
	MongoDBDefaultDataBase = "docker"
)

type MongDBClient struct {
	mongoDBReader *mongo.Client
	URI           string
	DataBase      string
}

func (m *MongDBClient) init() {
	// Handle Terminate Signal(Ctrl + C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = m.disconnectMongoDBClient()
	}()
}
func (m *MongDBClient) InsertNodeOne(node swarm.Node) (bool, error) {
	if err := m.RefreshMongoDBConnection(); err != nil {
		return false, err
	}
	collection := m.mongoDBReader.Database(m.DataBase).Collection("node")
	_, err := collection.InsertOne(context.Background(), node)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *MongDBClient) InsertServiceOne(service swarm.Service) (bool, error) {
	if err := m.RefreshMongoDBConnection(); err != nil {
		return false, err
	}
	collection := m.mongoDBReader.Database(m.DataBase).Collection("service")
	_, err := collection.InsertOne(context.Background(), service)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *MongDBClient) InsertTaskOne(task swarm.Task) (bool, error) {
	if err := m.RefreshMongoDBConnection(); err != nil {
		return false, err
	}
	collection := m.mongoDBReader.Database(m.DataBase).Collection("task")
	_, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *MongDBClient) UpdateNodeOne(node swarm.Node) (bool, error) {
	if err := m.RefreshMongoDBConnection(); err != nil {
		return false, err
	}
	updateOption := options.UpdateOptions{}
	updateOption.SetUpsert(true)
	collection := m.mongoDBReader.Database(m.DataBase).Collection("node")
	_, err := collection.UpdateOne(context.Background(), bson.D{{"id", node.ID}}, bson.D{{"$set", node}}, &updateOption)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *MongDBClient) UpdateServiceOne(service swarm.Service) (bool, error) {
	if err := m.RefreshMongoDBConnection(); err != nil {
		return false, err
	}
	updateOption := options.UpdateOptions{}
	updateOption.SetUpsert(true)
	collection := m.mongoDBReader.Database(m.DataBase).Collection("service")
	_, err := collection.UpdateOne(context.Background(), bson.D{{"id", service.ID}}, bson.M{"$set": bson.M{"dockerservice": service}}, &updateOption)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *MongDBClient) UpdatePrizesServiceOne(service service.PrizesService) (bool, error) {
	if err := m.RefreshMongoDBConnection(); err != nil {
		return false, err
	}
	updateOption := options.UpdateOptions{}
	updateOption.SetUpsert(true)
	collection := m.mongoDBReader.Database(m.DataBase).Collection("service")
	_, err := collection.UpdateOne(context.Background(), bson.D{{"id", service.DockerSerivce.ID}}, bson.D{{"$set", service}}, &updateOption)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *MongDBClient) UpdateTaskOne(task swarm.Task) (bool, error) {
	if err := m.RefreshMongoDBConnection(); err != nil {
		return false, err
	}
	updateOption := options.UpdateOptions{}
	updateOption.SetUpsert(true)
	collection := m.mongoDBReader.Database(m.DataBase).Collection("task")
	_, err := collection.UpdateOne(context.Background(), bson.D{{"id", task.ID}}, bson.D{{"$set", task}}, &updateOption)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *MongDBClient) FindNodeOne(NodeID string) (*swarm.Node, error) {
	var node swarm.Node
	if err := m.RefreshMongoDBConnection(); err != nil {
		return nil, err
	}
	collection := m.mongoDBReader.Database(m.DataBase).Collection("node")

	if err := collection.FindOne(context.Background(), bson.D{{"id", NodeID}}).Decode(&node); err != nil {
		return nil, err
	}
	return &node, nil
}

func (m *MongDBClient) FindServiceOne(serviceID string) (*swarm.Service, error) {
	var service swarm.Service
	if err := m.RefreshMongoDBConnection(); err != nil {
		return nil, err
	}
	collection := m.mongoDBReader.Database(m.DataBase).Collection("service")

	if err := collection.FindOne(context.Background(), bson.D{{"id", serviceID}}).Decode(&service); err != nil {
		return nil, err
	}
	return &service, nil
}

func (m *MongDBClient) FindPrizesServiceOne(serviceID string) (*service.PrizesService, error) {
	var prizeService service.PrizesService
	if err := m.RefreshMongoDBConnection(); err != nil {
		return nil, err
	}
	collection := m.mongoDBReader.Database(m.DataBase).Collection("service")

	if err := collection.FindOne(context.Background(), bson.D{{"id", serviceID}}).Decode(&prizeService); err != nil {
		return nil, err
	}
	return &prizeService, nil
}

func (m *MongDBClient) FindPrizesServiceFromPubkey(pubkey string) (*[]service.PrizesService, error) {
	fmt.Println(pubkey)
	var servicelist []service.PrizesService
	if err := m.RefreshMongoDBConnection(); err != nil {
		return nil, err
	}
	collection := m.mongoDBReader.Database(m.DataBase).Collection("service")

	cursor, err := collection.Find(context.Background(), bson.M{"createspec.pubkey": bson.M{"$eq": pubkey}}, nil)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var service service.PrizesService
		if err := cursor.Decode(&service); err != nil {
			return nil, err
		}
		servicelist = append(servicelist, service)
	}
	return &servicelist, nil
}

func (m *MongDBClient) FindTaskList(serviceID string) (*[]swarm.Task, error) {
	var tasklist []swarm.Task
	if err := m.RefreshMongoDBConnection(); err != nil {
		return nil, err
	}
	collection := m.mongoDBReader.Database(m.DataBase).Collection("task")

	cursor, err := collection.Find(context.Background(), bson.D{{"serviceid", serviceID}}, nil)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var task swarm.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasklist = append(tasklist, task)
	}
	return &tasklist, nil
}

// dialMongoDB connects a client to MongoDB server.
// Returns a MongoDB Client or any dialing error.
func (m *MongDBClient) dialMongoDB() error {
	if strings.TrimSpace(m.URI) == "" {
		return errors.New("uri format error")
	}
	var err error
	m.mongoDBReader, err = mongo.Connect(context.Background(), options.Client().ApplyURI(m.URI))
	if err != nil {
		return err
	}
	if err := m.mongoDBReader.Ping(context.TODO(), nil); err != nil {
		return err
	}
	return nil
}

// disconnectMongoDBClient disconnects a client from MongoDB server.
// Returns if there is any disconnection error.
func (m *MongDBClient) disconnectMongoDBClient() error {
	if m.mongoDBReader == nil {
		return errors.New("ErrNilMongoDBClient")
	}
	return m.mongoDBReader.Disconnect(context.Background())
}

// RefreshMongoDBConnection refreshes a client's connection with MongoDB server.
// Returns if there is any connection error.
func (m *MongDBClient) RefreshMongoDBConnection() error {
	if m.mongoDBReader == nil {
		err := m.dialMongoDB()
		if err != nil {
			return err
		}
		return nil
	}
	if err := m.mongoDBReader.Ping(context.TODO(), nil); err != nil {
		if err := m.disconnectMongoDBClient(); err != nil {
			m.mongoDBReader = nil
			return err
		}
		if err := m.dialMongoDB(); err != nil {
			m.mongoDBReader = nil
			return err
		}
	}
	return nil
}
