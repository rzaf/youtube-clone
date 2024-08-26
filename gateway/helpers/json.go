package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func WriteJsonError(w http.ResponseWriter, err string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	e := json.NewEncoder(w).Encode(map[string]string{"error": err})
	if e != nil {
		panic("JSON encoding failed. err:" + e.Error())
	}
}

func WriteJsonMessage(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	e := json.NewEncoder(w).Encode(map[string]string{"message": message})
	if e != nil {
		panic("JSON encoding failed. err:" + e.Error())
	}
}

func WriteJson(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic("JSON encoding failed. err:" + err.Error())
	}
}

func WriteEmpty(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
}

func WriteProtoJson(w http.ResponseWriter, data proto.Message, emitDefaultValues bool, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	m := protojson.MarshalOptions{EmitDefaultValues: emitDefaultValues}
	res, err := m.Marshal(data)
	if err != nil {
		panic("protojson.Marshal failed. err:" + err.Error())
	}
	_, err = w.Write(res)
	if err != nil {
		panic("Failed to write response to ResponseWriter!")
	}
}

func ReadJson(r *http.Request, data any) {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("JSON decoding failed. err:", err)
		panic(ServerError{Status: 400, Message: "request body should be a valid JSON!"})
	}
}

// func ReadJson(r *http.Request, data any) error {
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }