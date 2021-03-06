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

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.Handle("/pair-device",
		PairDeviceHandler(NewCreatePairDevice(db))).
		Methods(http.MethodPost)

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

func PairDeviceHandler(device Device) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("readall error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("raw request: %s\n", string(b))
		defer r.Body.Close()

		var p Pair
		if err := json.Unmarshal(b, &p); err != nil {
			log.Printf("unmarshal error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("raw request in struct: %#v\n", p)

		if err = device.Pair(p); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		fmt.Println("insert document success.")

		w.Write([]byte(`{"status":"active"}`))
	}
}

// Example from: net/http/server.go
// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

type Device interface {
	Pair(p Pair) error
}

type CreatePairDeviceFunc func(p Pair) error

func (fn CreatePairDeviceFunc) Pair(p Pair) error {
	return fn(p)
}

func NewCreatePairDevice(db *sql.DB) CreatePairDeviceFunc {
	return func(p Pair) error {
		// insert Pair obj
		_, err := db.Exec(`INSERT INTO pairs VALUES ($1,$2);`,
			p.DeviceID, p.UserID)
		if err != nil {
			log.Printf("can't insert doc %#v into table, error: %#v\n", p, err)
			return err
		}

		return nil
	}
}
