package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJsonError(w http.ResponseWriter, err string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	e := json.NewEncoder(w).Encode(map[string]string{"error": err})
	if e != nil {
		log.Fatal("JSON encoding failed")
	}
}

func WriteJsonMessage(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	e := json.NewEncoder(w).Encode(map[string]string{"message": message})
	if e != nil {
		log.Fatal("JSON encoding failed")
	}
}

func WriteJson(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Panic(err.Error())
	}
}

func ReadJson(r *http.Request, data any) {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(ServerError{Status: 400, Message: "invalid JSON request body, error:"})
	}
}

// func ReadJson(r *http.Request, data any) error {
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
