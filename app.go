package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	. "restapi-movies/config"
	. "restapi-movies/dao"
	. "restapi-movies/models"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var config = Config{}
var dao = MoviesDAO{}

// GetSamplePNG returns the example png
func GetSamplePNG(w http.ResponseWriter, r *http.Request) {
	respondWithPNG(w, http.StatusOK, nil)
}

// AllMoviesEndPoint fetches all movies
func AllMoviesEndPoint(w http.ResponseWriter, r *http.Request) {
	movies, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, movies)
}

//FindMovieEndPoint finds a movie by id
func FindMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, movie)
}

//CreateMovieEndPoint creates a movied record accordingly
func CreateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	movie.ID = bson.NewObjectId()
	if err := dao.Insert(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, movie)
}

//UpdateMovieEndPoint updates a movie record as requested
func UpdateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := dao.Update(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

//DeleteMovieEndPoint deletes a movie record per requested
func DeleteMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(movie); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithPNG(w http.ResponseWriter, code int, payload interface{}) {
	file, err := os.OpenFile("/Users/bruce/Downloads/coffee.png", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var b []byte
	if b, err = ioutil.ReadAll(file); err != nil {
		log.Fatal(err)
	}

	log.Printf("The length of image %d", len(b))
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Lenght", strconv.Itoa(len(b)))
	w.WriteHeader(code)
	w.Write(b)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func test(w http.ResponseWriter, code int, msg string) {
	fmt.Println("what the hell!")
}

/*
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}
*/

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/movies", AllMoviesEndPoint).Methods("GET")
	r.HandleFunc("/movies", CreateMovieEndPoint).Methods("POST")
	r.HandleFunc("/movies", UpdateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/movies", DeleteMovieEndPoint).Methods("DELETE")
	r.HandleFunc("/movies/{id}", FindMovieEndPoint).Methods("GET")
	r.HandleFunc("/png", GetSamplePNG).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
