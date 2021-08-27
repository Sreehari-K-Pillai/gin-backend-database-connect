package main

import (
	// "time"
	"net/http"
	"github.com/gin-gonic/gin"
	// "github.com/mattn/go-colorable"
	"github.com/gin-contrib/cors"

	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    //"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson"

	"context"
	"log"
	"time"
	//"fmt"
)

type album struct {
    // User     string `form:"user" json:"user" xml:"user" `
	// Password string `form:"password" json:"password" xml:"password"`
    // ID     string  `form: "id" json:"id"`
    Title  string  `form: "title" json:"title"`
    EpisodeNo string  `form: "episodeNo" json:"episodeNo"`
    // Price  float64 `form: "price" json:"price"`
}

// var albums = []album{
//     {User: "Sreehari", Password: "123456", ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
//     {User: "Gunjan", Password: "12345", ID: "2", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
//     {User: "Tanmay", Password: "1234", ID: "3", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
//     {User: "Omkar", Password: "123", ID: "4", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }


func main() {
	r := gin.Default()
	// gin.DefaultWriter = colorable.NewColorableStdout()

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return origin == "*"
	// 	},
	// 	MaxAge: 12 * time.Hour,
	// }))
		r.Use(cors.Default())
	// emp := r.Group("/auth")
	// emp.POST("/registration", Register)
	auth := r.Group("/data")
	auth.GET("/getuser", getUser)
	auth.POST("/adduser", addUser)
	r.Run(":3030")
}

func addUser(c *gin.Context) {

	client,err:= mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx,_:= context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)


    var newAlbum *album

    if err := c.BindJSON(&newAlbum); err != nil {
		log.Fatal(err)
        
    }

	_,err = client.Database("quickstart").Collection("episode").InsertOne(ctx,newAlbum)

	if err != nil {
		log.Fatal(err)
	}

    c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getUser(c *gin.Context) {

	client,err:= mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx,_:= context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	// databases,err := client.ListDatabaseNames(ctx,bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	cursor,err := client.Database("quickstart").Collection("episode").Find(ctx,bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var episodes []bson.M
	if err := cursor.All(ctx,&episodes); err!= nil {
		log.Fatal(err)
	}

    c.IndentedJSON(http.StatusOK, episodes)
}

