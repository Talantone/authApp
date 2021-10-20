package apiserver

import (
	"authApp"
	"authApp/internal/app/store"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("starting api server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil

}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}

func (s *APIServer) configureRouter() {
	//s.router.HandleFunc("/auth", s.Auth())
	s.router.HandleFunc("/reg", s.Reg()).Methods("POST")
	s.router.HandleFunc("/people", s.GetPeople()).Methods("GET")
	s.router.HandleFunc("/user/{id}", s.GetPersonByID()).Methods("GET")

}

func (s *APIServer) Reg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		var user authApp.CreateUser
		json.NewDecoder(r.Body).Decode(&user)
		collection := s.store.Client.Database("authApp").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), time.Second*15)
		result, _ := collection.InsertOne(ctx, user)
		json.NewEncoder(w).Encode(result)
	}
}
func (s *APIServer) GetPeople() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		var people []authApp.User
		collection := s.store.Client.Database("authApp").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), time.Second*15)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var user authApp.User
			cursor.Decode(&user)
			people = append(people, user)
		}
		if err := cursor.Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}

		json.NewEncoder(w).Encode(people)

	}
}
func (s *APIServer) GetPersonByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])
		var user authApp.User

		collection := s.store.Client.Database("authApp").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), time.Second*15)
		err := collection.FindOne(ctx, authApp.User{ID: id}).Decode(&user)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}

		fmt.Print(user)
		json.NewEncoder(w).Encode(user)
	}
}
