package router

import (
	controller "mongoapi/Controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/allmovies", controller.GetAllRecords).Methods("GET")

	router.HandleFunc("/api/deleteall", controller.DeleteAll).Methods("DELETE")
	router.HandleFunc("/api/deletemovie/{id}", controller.DeleteMovie).Methods("DELETE")

	router.HandleFunc("/api/movie", controller.CreateMovie).Methods("POST")

	router.HandleFunc("/api/watchedmovie/{id}", controller.MarkMovieAsWatched).Methods("PUT")

	return router
}
