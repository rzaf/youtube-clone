package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

func readJson(r *http.Request, data any) {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("JSON decoding failed. err:", err)
		panic(ServerError{Status: 400, Message: "request body should be a valid JSON!"})
	}
}

func ParseReq(r *http.Request, data map[string]any) {
	if data == nil {
		panic("`data` argument is not initialized")
	}
	contentType := r.Header.Get("Content-Type")
	if r.Method == "GET" {
		if contentType != "" {
			panic(ServerError{Status: 400, Message: "request body not supported on GET reuqests"})
		}
		for k, v := range r.URL.Query() {
			data[k] = v[0]
		}
		return
	}
	contentLength, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	if contentLength == 0 {
		if contentType == "" {
			return // giving empty body
		}
		panic(ServerError{Status: 400, Message: "empty request body"})
	}
	switch contentType {
	case "application/json":
		readJson(r, &data)
	case "application/x-www-form-urlencoded":
		r.ParseForm()
		for k, v := range r.PostForm {
			data[k] = v[0]
		}
		fmt.Printf("PostForm: %v\n", r.PostForm)
	default:
		if strings.Contains(contentType, "multipart/form-data") {
			err := r.ParseMultipartForm(5000)
			if err != nil {
				panic(ServerError{Status: 400, Message: "failed to read request body"})
			}
			for k, v := range r.MultipartForm.Value {
				data[k] = v[0]
			}
			fmt.Printf("MultipartForm.File : %v\n", r.MultipartForm.File)
			fmt.Printf("MultipartForm.Value: %v\n", r.MultipartForm.Value)
			break
		}
		panic(ServerError{Status: 400, Message: "invalid request content-ype"})
	}
}

// func ReadJson(r *http.Request, data any) error {
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
