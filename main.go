package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/guillermo-st/locus/chat"
)

func main() {

	port := "8000"
	l := log.New(os.Stdout, "locus", log.LstdFlags)
	r := mux.NewRouter()

	//Handle routes
	chatHandler := chat.NewChat(l)
	r.Handle("/chat/{roomName}", chatHandler)

	http.ListenAndServe(":"+port, r)
}
