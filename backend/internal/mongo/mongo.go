package mongo

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var ctx context.Context
var cancel context.CancelFunc
var dataBase string

func connect(uri string)(*mongo.Client, context.Context, context.CancelFunc, error) {
    // Set credentials
    credential := options.Credential{
        Username: os.Getenv("MONGO_USERNAME"),
        Password: os.Getenv("MONGO_PASSWORD"),
    }

    // Set mongo context
    ctx, cancel := context.WithCancel(context.Background())

    // Start connection to mongo
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(credential))

    // Create index for field username in collection users
    collection := client.Database(dataBase).Collection("users")

    _, errIndex := collection.Indexes().CreateOne(
        ctx, 
        mongo.IndexModel{
            Keys:    bson.D{{Key: "username", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    )
    if errIndex != nil {
		log.Printf("Error: Create index of users")
		panic(err)
	}

    // Create index for field userId in collection santanderTokens
    collection = client.Database(dataBase).Collection("santanderTokens")

    _, errIndex = collection.Indexes().CreateOne(
        ctx, 
        mongo.IndexModel{
            Keys:    bson.D{{Key: "userId", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    )
    if errIndex != nil {
		log.Printf("Error: Create index of santanderTokens")
		panic(err)
	}


    return client, ctx, cancel, err
}

func Close(){
    log.Printf("Closing mongo connection")
    defer cancel()
     
    defer func() {
        if err := client.Disconnect(ctx); err != nil {
            log.Printf("Error: Discconecting mongo")
            panic(err)
        }
    }()
}
 
func InsertOne(col string, doc interface{})(*mongo.InsertOneResult, error) {
    log.Printf("Mongo: Insert one")
    
    collection := client.Database(dataBase).Collection(col)
    
    return collection.InsertOne(ctx, doc)
}

func InsertMany(col string, docs []interface{})(*mongo.InsertManyResult, error) {
    log.Printf("Mongo: Insert many")

    collection := client.Database(dataBase).Collection(col)
     
    return collection.InsertMany(ctx, docs)
}

func Query(col string, query, field interface{}) (result *mongo.Cursor, err error) {
    log.Printf("Mongo: Query")

    collection := client.Database(dataBase).Collection(col)

    return collection.Find(ctx, query, options.Find().SetProjection(field))
}

func Aggregate(col string, pipeline interface{}) (result *mongo.Cursor, err error) {
    log.Printf("Mongo: Aggregate")

    collection := client.Database(dataBase).Collection(col)

    return collection.Aggregate(ctx, pipeline)
}

func UpdateOne(col string, filter, update interface{}, opts *options.UpdateOptions) (result *mongo.UpdateResult, err error) {
    log.Printf("Mongo: Update one")

    collection := client.Database(dataBase).Collection(col)

    return collection.UpdateOne(ctx, filter, update, opts)
}

func UpdateMany(col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {
    log.Printf("Mongo: Update many")

    collection := client.Database(dataBase).Collection(col)

    return collection.UpdateMany(ctx, filter, update)
}

func ReplaceOne(col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {
    log.Printf("Mongo: Replace one")

    collection := client.Database(dataBase).Collection(col)

    return collection.ReplaceOne(ctx, filter, update)
}

func DeleteOne(col string, query interface{}) (result *mongo.DeleteResult, err error) {
    log.Printf("Mongo: Delete one")

    collection := client.Database(dataBase).Collection(col)

    return collection.DeleteOne(ctx, query)
}

func DeleteMany(col string, query interface{}) (result *mongo.DeleteResult, err error) {
    log.Printf("Mongo: Delete many")

    collection := client.Database(dataBase).Collection(col)

    return collection.DeleteMany(ctx, query)
}

func GetCtx() (context.Context) {
    log.Printf("Get mongo context")
    return ctx
}
 
func Start() {
    var err error

    log.Printf("Starting mongo connection")

    dataBase = os.Getenv("MONGO_DATABASE")

    client, ctx, cancel, err = connect(os.Getenv("MONGO_PATH"))
    if err != nil {
        log.Printf("Error: Mongo connection")
        panic(err)
    }
}