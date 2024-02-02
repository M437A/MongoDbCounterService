package connections

import (
	"context"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoPool *mongo.Client
	once      sync.Once
	ctx       context.Context
	cancel    context.CancelFunc
)

func GetMongoClient() *mongo.Client {
	return mongoPool
}

func GetMongoDB(dataBaseName string) *mongo.Database {
	return mongoPool.Database(dataBaseName)
}

func connectToMongoDB(ctx context.Context) {
	once.Do(func() {
		ctx, cancel = context.WithCancel(ctx)

		clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL")).SetAuth(options.Credential{
			Username: os.Getenv("MONGO_USERNAME"),
			Password: os.Getenv("MONGO_PASSWORD"),
		})

		var err error
		mongoPool, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Error initializing MongoDB connection pool: %v", err)
		}

		if err := ping(ctx, mongoPool); err != nil {
			log.Fatalf("Error pinging MongoDB: %v", err)
		}
		log.Printf("MongoDB was ininzilizate secssesful")
	})
}

func ping(ctx context.Context, client *mongo.Client) error {
	err := client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Printf("Ping to mongoDB was secssesful")

	return nil
}

func CloseMongoDBConnection() {
	cancel()
	if mongoPool != nil {
		_ = mongoPool.Disconnect(ctx)
		log.Printf("MongoDB cloused secssesful")
	}
}
