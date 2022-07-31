package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	model "mongoapi/Model"
	"net/http"

	//"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://sakshi1104:14371818@cluster0.judp1.mongodb.net/?retryWrites=true&w=majority"
const dbName = "primevideo"
const collectionName = "watchlist"

//original reference of mongoDB collection
var collection *mongo.Collection

//connect with mongodb
//init function only runs once, where the program is first run
// @Title Initialises a client
// @Description Intialises a client to connect with mongodb
func init() {
	//client options
	clientOption := options.Client().ApplyURI(connectionString)

	//client used to connect to mongodb (connection request made)
	//context -- whenever connecting to an outside machine, or online server - how long connection is made for, what happens when connection goes off etc
	//TODO(is the type of context) is best option when you dont know which type of context to use
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection successfully made")

	collection = client.Database(dbName).Collection(collectionName)

	//collection instance
	fmt.Println("Collection instance is ready")
}

//insert 1 record - movie
//inserted i lowercase as helper method will not be exported
// @Title Insert One Movie
// @Description Helper method that inserts one movie into the database and gives it an ID
// @Param movie body model.PrimeVideo true "info about the movie"
func insertOneMovie(movie model.PrimeVideo) {
	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted 1 movie in DB with id: ", inserted.InsertedID)
}

//update 1 record - movie
// @Title Update One Movie
// @Description Helper method that updates watched status to true according to the provided movie id
// @Param movieId body string true "movie id to be updated"
func updateOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified movie count: ", result.ModifiedCount)
}

//delete 1 record - movie
// @Title Deletes One Movie
// @Description Helper method that deletes one movie from the database according to the provided movie id
// @Param movieId body string true "movie id to be deleted"
func deleteOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}

	//deleteCount is actually an integer value unlike inserted
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Movie got deleted with delete count: ", deleteCount)
}

//delete all records from mongodb
// @Title Deletes All Movie
// @Description Helper method that deletes all movies from the database and returns a delete count
func deleteAll() int64 {
	//filter := bson.D{{}}
	//empty paranthesis means no filter select all
	//deleteResult is key-value pair of int64 type
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of movies deleted: ", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

//get all from mongodb
// @Title Get All Movies
// @Description Helper method that gets all the movies/records in the database and returns them
func getAll() []primitive.M {
	//cursor is a type of object adn you have to loop through it to obtain the data
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var records []primitive.M

	for cursor.Next(context.Background()) {
		var record bson.M
		err := cursor.Decode(&record)
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, record)
	}

	defer cursor.Close(context.Background())
	return records
}

//Actual controller that we are using
//getting movies
// @Title Get All Records
// @Description Actual controller that gets all the movies/records in the database and returns them
// @Param respone body  true "creates a writer for json"
// @Param request body pointer true "creates a reader for json"
// @Router /api/allmovies [get]
func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	//or application/x-www-form-urlencode
	w.Header().Set("Content-Type", "apllication/json")
	allRecords := getAll()
	json.NewEncoder(w).Encode(allRecords)
}

// creating movie
// @Title Create a Records
// @Description Actual controller that creates a movie
// @Param respone body  true "creates a writer for json"
// @Param request body pointer true "creates a reader for json"
// @Router /api/movie [post]
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "apllication/json")
	w.Header().Set("Allow-Control_Allow-Methods", "POST")

	var movie model.PrimeVideo
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

// marking the movie as watched
// @Title Marking movie as watched
// @Description Actual controller that marks a movie as watched
// @Param respone body  true "creates a writer for json"
// @Param request body pointer true "creates a reader for json"
// @Router /api/watchedmovie/{id} [put]
func MarkMovieAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "apllication/json")
	w.Header().Set("Allow-Control_Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

// deleting 1 movie
// @Title Deleting a movie
// @Description Actual controller that deletes a movie
// @Param respone body  true "creates a writer for json"
// @Param request body pointer true "creates a reader for json"
// @Router /api/deletemovie/{id} [delete]
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "apllication/json")
	w.Header().Set("Allow-Control_Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

//deleting all
// @Title Deleting all movies
// @Description Actual controller that deletes all movies
// @Param respone body  true "creates a writer for json"
// @Param request body pointer true "creates a reader for json"
// @Router /api/deleteall [delete]
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "apllication/json")
	w.Header().Set("Allow-Control_Allow-Methods", "DELETE")

	count := deleteAll()
	json.NewEncoder(w).Encode(count)
}
