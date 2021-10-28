package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/julioisaac/daxxer-api/cmd"
)

type muxRouter struct{}

var (
	muxDispatcher = mux.NewRouter()
)

func NewMuxRouter() cmd.Router {
	return &muxRouter{}
}

func (m *muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("GET")
}

func (m *muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("POST")
}

func (m *muxRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("PUT")
}

func (m *muxRouter) DELETE(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("DELETE")
}

func (m *muxRouter) SERVE(port string) {
	err := http.ListenAndServe(port, muxDispatcher)
	if err != nil {
		fmt.Printf("Error trying running Mux")
		return
	}
}
