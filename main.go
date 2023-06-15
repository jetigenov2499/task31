package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %s", err)
	}
	log.Printf("Connected to Redis: %s", pong)
}

func getCachedResponse(name string) (string, error) {
	cachedResponse, err := redisClient.Get(redisClient.Context(), name).Result()
	if err != nil {
		return "", err
	}
	return cachedResponse, nil
}

func cacheResponse(name, response string) error {
	err := redisClient.Set(redisClient.Context(), name, response, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	cachedResponse, err := getCachedResponse(name)
	if err == nil && cachedResponse != "" {
		fmt.Fprintln(w, cachedResponse)
		return
	}

	response := fmt.Sprintf("Hello, %s!", name)
	fmt.Fprintln(w, response)

	err = cacheResponse(name, response)
	if err != nil {
		log.Printf("Failed to cache response: %s", err)
	}
}

func main() {
	initRedis()

	http.HandleFunc("/api/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
