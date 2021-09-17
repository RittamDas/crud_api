package main

import (
	"context"
	"crud-api/crud"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	var err error
	crud.Client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://rittam:abcd1234@cluster0.whh2m.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = crud.Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer crud.Client.Disconnect(ctx)
	err = crud.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")
	router := mux.NewRouter()
	router.HandleFunc("/users", crud.CreateUser).Methods("POST")
	router.HandleFunc("/users/{_id}", crud.GetUserById).Methods("GET")
	router.HandleFunc("/users", crud.GetUsers).Methods("GET")
	router.HandleFunc("/users/{_id}", crud.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{_id}", crud.GetUserById).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}
