package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Pair struct {
	DeviceID int64 `json:"DeviceID"`
	UserID   int64 `json:"UserID"`
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

	// open database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("connect to database error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// insert Pair obj
	_, err = db.Exec(`INSERT INTO pairs VALUES ($1,$2);`,
		rawRequest.DeviceID, rawRequest.UserID)
	if err != nil {
		log.Printf("can't insert doc %#v into table, error: %#v\n", rawRequest, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("insert document success.")

	w.Write([]byte(`{"status":"active"}`))
}
