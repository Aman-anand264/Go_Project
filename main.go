package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gofr.dev/pkg/gofr"
)

type Book struct {
	ISBN   string `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	mongoURI    = "mongodb://localhost:27017"
	database    = "bookstore"
	collection  = "books"
	mongoClient *mongo.Client
)

func initMongoDB() error {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}
	fmt.Println("Connected to MongoDB!")

	mongoClient = client
	return nil
}

func main() {
	err := initMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	app := gofr.New()

	app.POST("/book/add", func(ctx *gofr.Context) (interface{}, error) {
		var book Book
		err := json.NewDecoder(ctx.Request().Body).Decode(&book)
		if err != nil {
			return nil, err
		}

		log.Printf("ISBN: %s", book.ISBN)

		booksCollection := mongoClient.Database(database).Collection(collection)
		_, err = booksCollection.InsertOne(context.Background(), book)
		if err != nil {
			return nil, err
		}
		return "Book added to the bookstore successfully", nil
	})

	app.GET("/books/list", func(ctx *gofr.Context) (interface{}, error) {
		var books []Book

		booksCollection := mongoClient.Database(database).Collection(collection)
		cursor, err := booksCollection.Find(context.Background(), bson.M{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			var book Book
			if err := cursor.Decode(&book); err != nil {
				return nil, err
			}
			books = append(books, book)
		}
		if err := cursor.Err(); err != nil {
			return nil, err
		}
		return books, nil
	})

	app.PUT("/books/list/{isbn}", func(ctx *gofr.Context) (interface{}, error) {

		rawPath := ctx.Request().URL.Path

		pathParts := strings.Split(rawPath, "/")
		isbn := pathParts[3]

		fmt.Println("Handling PUT request for ISBN:", isbn)
		var updatedBook Book
		err := json.NewDecoder(ctx.Request().Body).Decode(&updatedBook)
		if err != nil {
			return nil, err
		}
		booksCollection := mongoClient.Database(database).Collection(collection)
		filter := bson.M{"isbn": isbn}
		update := bson.M{"$set": updatedBook}

		result, err := booksCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return nil, err
		}
		if result.ModifiedCount == 0 {
			return nil, fmt.Errorf("No book found with ISBN: %s", isbn)
		}

		return "Book updated successfully", nil
	})

	app.DELETE("/book/remove/{isbn}", func(ctx *gofr.Context) (interface{}, error) {
		rawPath := ctx.Request().URL.Path
		pathParts := strings.Split(rawPath, "/")
		isbn := pathParts[3]
		booksCollection := mongoClient.Database(database).Collection(collection)
		result, err := booksCollection.DeleteOne(context.Background(), bson.M{"isbn": isbn})
		if err != nil {
			return nil, err
		}
		if result.DeletedCount == 0 {
			return nil, fmt.Errorf("No books found with ISBN: %s", isbn)
		}
		return "Book with ISBN removed successfully", nil
	})

	app.Start()
}
