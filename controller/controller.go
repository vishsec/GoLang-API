package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"vishsec.dev/goapi/model"
)

const dbName = "netflix"
const colName = "watchlist"                                                                        //connection name

var collection *mongo.Collection                                                                   //helps you take reference of mongodb collection

func init() {

	err := godotenv.Load()
	if err!=nil {
		log.Fatalf(".env cant be loaded %v", err)
	}

	connectionString := os.Getenv("API_KEY")
	if connectionString == "" {
		log.Fatal("API_KEY not set in .env file")
	}


	clientOption := options.Client().ApplyURI(connectionString)                                    // basic syntax for initializing a connection and configure mongodb client in go
	                                                                                               // firing up connnection request
	client, err := mongo.Connect(context.TODO(), clientOption)                                     // (basically keeps the connection alive when there is something todo)context has deadlines and cancellation signals whenever we make contact with the server
	                                                                                               //types of context https://pkg.go.dev/context#Background
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connected")

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("collection instance is ready")

}

// mongo helper methods

func insertOneMovie(movie model.Netflix){
	inserted, err := collection.InsertOne(context.Background(), movie) 
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted movie id: " , inserted.InsertedID)
}

func updateOneMovie(movieId string){                                  
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id} //bson map to query the id in db
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update) //result gives you count of updated docs
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)


}

func deleteOneMovie(movieId string){
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	// update := bson.M{"$set": bson.M{"watched": true}}
	delCount, err := collection.DeleteOne(context.Background(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println( "delete count: ", delCount.DeletedCount)
}

func deleteAll() int64{
	filter := bson.D{{}}
	count, err := collection.DeleteMany(context.Background(), filter, nil)
	if(err!=nil){
		log.Fatal(err)
	}

	fmt.Println("delete counts : ", count.DeletedCount)
	return count.DeletedCount
}

func getAll()  []primitive.M{
	cur , err := collection.Find(context.Background() ,bson.D{})
	if(err!=nil){
		log.Fatal(err)
	}

	var movies []primitive.M

	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		if(err!=nil){
			log.Fatal(err)
		}

		movies = append(movies, movie)
	}

	defer cur.Close(context.Background())
	return movies

}


// Actual controller methods 

func GetAllMovies(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","x-www-form-urlencode")
	allMovies := getAll()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "x-www-form-urlencode")
	w.Header().Set("Allow-Control_Allow_Methods", "POST")

	var movie model.Netflix

	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkMovie(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control_Allow_Methods", "POST")

	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteMovie(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control_Allow_Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])

	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAll(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control_Allow_Methods", "DELETE")

	count := deleteAll()
	json.NewEncoder(w).Encode(count)
}
