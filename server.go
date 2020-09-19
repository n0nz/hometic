package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Pair struct {
	DeviceID int64 `json:"device_id"`
	UserID   int64 `json:"user_id"`
}

func main() {
	fmt.Println("hello hometic : I'm Gopher!!")

	r := mux.NewRouter()
	r.HandleFunc("/pair-device", PairDeviceHandler).Methods(http.MethodPost)

	r.Use(Middleware)

	sAddr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	server := http.Server{
		Addr:    sAddr,
		Handler: r,
	}

	log.Printf("starting at %s\n", sAddr)
	log.Fatal(server.ListenAndServe())
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("middleware", r.URL)
		h.ServeHTTP(w, r)
	})
}

func PairDeviceHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("readall error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("raw request: %s\n", string(b))
	defer r.Body.Close()

	var rawRequest Pair
	if err := json.Unmarshal(b, &rawRequest); err != nil {
		log.Printf("unmarshal error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("raw request in struct: %#v\n", rawRequest)

	w.Write([]byte(`{"status":"active"}`))
}
