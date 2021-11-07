package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/guillermo-st/locus/chat"
	"github.com/guillermo-st/locus/handlers"
)

func main() {

	port := "8000"
	l := log.New(os.Stdout, "locus", log.LstdFlags)
	r := mux.NewRouter()

	var rooms = make(map[string]*chat.Room)

	//Handle routes
	rh := handlers.NewRooms(rooms, l)

	r.HandleFunc("/rooms", rh.GetRooms).Methods("GET")
	r.HandleFunc("/rooms", rh.CreateRoom).Methods("POST")

	http.ListenAndServe(":"+port, r)
}
