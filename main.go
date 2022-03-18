package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
)

func main() {
	//init router and redis
	router := mux.NewRouter()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	cacheKey := "mycache"

	//API 1
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"message": "Welcome to Dockerized app",
		}
		json.NewEncoder(rw).Encode(response)
		log.Println("Base url Hit on ", time.Now())
	})

	//API 2
	router.HandleFunc("/user/{name}", func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		var message string
		if name == "" {
			message = "Hello World"
		} else {
			message = "Hello " + name
		}
		response := map[string]string{
			"message": message,
		}
		json.NewEncoder(rw).Encode(response)
		log.Println("/{name} Hit on ", time.Now())

		// we can call set with a `Key` and a `Value`.
		err := client.HSet(cacheKey, name, "").Err()
		// if there has been an error setting the value
		// handle the error
		if err != nil {
			fmt.Println(err)
		}
	})

	//API 3
	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("Hello from Redis, %s", pong)

		response := map[string]string{
			"message": message,
		}
		json.NewEncoder(w).Encode(response)
	})

	//API 4
	router.HandleFunc("/names", func(w http.ResponseWriter, r *http.Request) {
		val, errGet := client.HGetAll(cacheKey).Result()
		if errGet != nil {
			fmt.Println(errGet)
		}
		message := fmt.Sprintf("All names stored: %s", val)

		response := map[string]string{
			"message": message,
		}
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Server is running!")
	fmt.Println(http.ListenAndServe(":8081", router))
}