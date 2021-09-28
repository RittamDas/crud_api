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
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Created   time.Time          `json:"created,omitempty" bson:"created,omitempty"`
	FirstName string             `json:"firstName,omitempty" bson:"firstName,omitempty" validate:"required"`
	LastName  string             `json:"lastName,omitempty" bson:"lastName,omitempty" validate:"required"`
	Age       struct {
		value    int64
		interval int64
	} `json:"age,omitempty" bson:"age,omitempty"`
	Mobile string `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Active bool   `json:"active" bson:"active"`
}

var Client *mongo.Client

func (u *User) Validate() error{
	validate := validator.New()
	return validate.Struct(u)
}

func CreateUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("content-type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	err := user.Validate()
	if err != nil{
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}
	user.Created = time.Now()
	user.Active = true
	collect := Client.Database("users").Collection("user_data")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, _ := collect.InsertOne(ctx, user)
	json.NewEncoder(rw).Encode(res)
	rw.Write([]byte(`{ "message": "` + user.FirstName + `" }`))
}

func GetUserById(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")
	p := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(p["id"])
	var user User
	collect := Client.Database("users").Collection("user_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collect.FindOne(ctx, bson.D{{"_id", id}, {"active", true}}).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{ "message": "` + user.ID.Hex() + `" }`))
		return
	}
	if !user.Active {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{ "message" : "User Not Found" }`))
		return
	}
	json.NewEncoder(rw).Encode(user)
}

func UpdateUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")
	p := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(p["id"])
	var user User
	collect := Client.Database("users").Collection("user_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collect.FindOne(ctx, bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if !user.Active {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{ "message" : "User Not Found" }`))
		return
	}
	_, err = collect.DeleteOne(ctx, User{ID: id})
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
	id, _ := primitive.ObjectIDFromHex(p["id"])
	var user User
	collect := Client.Database("users").Collection("user_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collect.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		rw.Write([]byte(id.Hex()))
		return
	}
	if !user.Active {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{ "message" : "User Not Found" }`))
		return
	}
	err = collect.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.D{{"$set", bson.D{{"active", false}}}}).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		rw.Write([]byte(id.Hex()))
		return
	}
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
		if user.Active {
			users = append(users, user)
		}
	}
	json.NewEncoder(rw).Encode(users)
}
