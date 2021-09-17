package crud

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	_id       primitive.ObjectID
	created   primitive.DateTime
	firstName string
	lastName  string
	age       struct {
		value    int64
		interval int64
	}
	mobile string
	active bool
}

var Client *mongo.Client

func CreateUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("content-type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	collect := Client.Database("users").Collection("user_data")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, _ := collect.InsertOne(ctx, user)
	json.NewEncoder(rw).Encode(res)
}

func GetUserById(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")
	p := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(p["_id"])
	var user User
	collect := Client.Database("users").Collection("user_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collect.FindOne(ctx, User{_id: id}).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if !user.active {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{ "message" : "User Not Found" }`))
		return
	}
	json.NewEncoder(rw).Encode(user)
}

func UpdateUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")
	p := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(p["_id"])
	var user User
	collect := Client.Database("users").Collection("user_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collect.FindOne(ctx, User{_id: id}).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if !user.active {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{ "message" : "User Not Found" }`))
		return
	}
	_, err = collect.DeleteOne(ctx, User{_id: id})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	var newUser User
	_ = json.NewDecoder(r.Body).Decode(&newUser)
	res, _ := collect.InsertOne(ctx, newUser)
	json.NewEncoder(rw).Encode(res)
}

func DeleteUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")
	p := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(p["_id"])
	var user User
	collect := Client.Database("users").Collection("user_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collect.FindOne(ctx, User{_id: id}).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if !user.active {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{ "message" : "User Not Found" }`))
		return
	}
	user.active = false
	json.NewEncoder(rw).Encode(user)
}

func GetUsers(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")
	var users []User
	collect := Client.Database("users").Collection("user_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collect.Find(ctx, bson.M{})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		if user.active {
			users = append(users, user)
		}
	}
	json.NewEncoder(rw).Encode(users)
}
