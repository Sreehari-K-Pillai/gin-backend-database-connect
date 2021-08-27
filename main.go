package main


import (
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson"

	"context"
	"log"
	"time"
	"fmt"
)

func main() {

	client,err:= mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx,_:= context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	

	err = client.Ping(ctx,readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	databases,err := client.ListDatabaseNames(ctx,bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases)

	defer client.Disconnect(ctx)

	//Create database and collections

	quickstartDatabase := client.Database("quickstart")
	podcastCollection := quickstartDatabase.Collection("podcast")
	episodeCollection := quickstartDatabase.Collection("episode")

	//Inser one object

	podcastResult,err := podcastCollection.InsertOne(ctx,bson.D{
		// {Key:"title",Value:"DBServer"},
		// {Key:"author",Value:"Mongo"},
		{"title","DBServer"},
		{"author","Mongo"},
		{"languages",bson.A{"C","C++","Python","JAVA"}},
	})

	if err != nil {
		log.Fatal(err)
	}


	fmt.Println(podcastResult.InsertedID)

	//Insert multiple objects

	episodeResult,err := episodeCollection.InsertMany(ctx,[]interface{}{
	bson.D{
		{"podcast",podcastResult.InsertedID},
		{"title","Episode"},
		{"episodeNo",1},
	},
	bson.D{
		{"podcast",podcastResult.InsertedID},
		{"title","Episode 10"},
		{"episodeNo",10},
	},
	})

	if err != nil {
		log.Fatal(err)
	}


	fmt.Println(episodeResult.InsertedIDs)

	// Getting data from database

	cursor,err := episodeCollection.Find(ctx,bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var episodes []bson.M
	if err := cursor.All(ctx,&episodes); err!= nil {
		log.Fatal(err)
	}

//	fmt.Println(episodes)

for _,e := range episodes {
	fmt.Println(e)
}

}