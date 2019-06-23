package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type Person struct {
	//ID        string   `json:"_id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City string `json:"city, omitempty"`
	State string `json:"state, omitempty"`
}

var collection *mongo.Collection
var people []Person

func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request)  {
	w.Write(make([]byte, 7))
	//log.Println(req.Method)
	fmt.Println(req.Method)
}

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request)  {
	//person := Person{}
	//_ = json.NewDecoder(req.URL.Query()["lastname"]).Decode(&person)
	res := collection.FindOne(context.TODO(), bson.D{{"firstname", req.URL.Query()["firstname"]}})
	log.Println()
	json.NewEncoder(w).Encode(res)
	return

}

func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request)  {
	//w.Header().Set("content-type", "application/json")
	person := Person{}
	_ = json.NewDecoder(req.Body).Decode(&person)
	collection.InsertOne(context.TODO(), person)
}

func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request)  {
	person := Person{}
	_ = json.NewDecoder(req.Body).Decode(&person)
	delete_result, err := collection.DeleteMany(context.TODO(), bson.D{{"firstname", person.Firstname}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v documents deleted from people collection\n", delete_result.DeletedCount)
}

func main() {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	defer client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//collection able to query on people collection of test_database
	collection = client.Database("test_database").Collection("people")




	router := mux.NewRouter()

	//people = append(people, Person{ID:"1", Firstname:"arif", Lastname:"ibrahim", Address:&Address{City:"Dhaka", State:"Dhaka"}})
	//people = append(people, Person{ID:"2", Firstname:"mamun", Lastname:"al", Address:&Address{City:"Dhaka", State:"Barishal"}})

	//router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people",DeletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}
