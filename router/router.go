package router

import (
	"github.com/gorilla/mux"
	"vishsec.dev/goapi/controller"
)

func Router() *mux.Router	{

	router := mux.NewRouter()
	router.HandleFunc("/api/movies", controller.GetAllMovies).Methods("GET")
	router.HandleFunc("/api/movie", controller.CreateMovie).Methods("POST")
	router.HandleFunc("/api/movie/{id}", controller.MarkMovie).Methods("POST")
	router.HandleFunc("/api/movie/{id}", controller.DeleteMovie).Methods("DELETE")
	router.HandleFunc("/api/deleteAll", controller.DeleteAll).Methods("DELETE")

	return router
}