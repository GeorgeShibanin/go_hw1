package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
	"twitterLikeHW/handlers"
	"twitterLikeHW/storage"
	"twitterLikeHW/storage/mongostorage"
)

func NewServer() *http.Server {
	r := mux.NewRouter()

	mongoUrl := os.Getenv("MONGO_URL")
	//uri := "mongodb://localhost:27017/twitterPosts"
	mongoStorage := mongostorage.NewStorage(mongoUrl)

	handler := &handlers.HTTPHandler{
		StorageOld: make(map[storage.PostId]storage.Post),
	}
	storageType := os.Getenv("STORAGE_MODE")
	if storageType == "inmemory" {
		handler = &handlers.HTTPHandler{
			StorageOld: make(map[storage.PostId]storage.Post),
		}
	} else {
		handler = &handlers.HTTPHandler{
			Storage: mongoStorage,
		}
	}

	r.HandleFunc("/", handlers.HandleRoot).Methods("GET", "POST")
	r.HandleFunc("/api/v1/posts", handler.HandleCreatePost).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/posts/{postId}", handler.HandleGetPosts).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/{userId}/posts", handler.HandleGetUserPosts).Methods(http.MethodGet)
	r.HandleFunc("/maintenance/ping", handlers.HandlePing).Methods(http.MethodGet)

	port := "8080"
	if value, ok := os.LookupEnv("SERVER_PORT"); ok {
		port = value
	}

	return &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func main() {
	srv := NewServer()
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
