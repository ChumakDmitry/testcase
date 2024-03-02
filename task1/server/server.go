package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"task1/internal"
)

type data struct {
	Result []int64
}

type APIServer struct {
	config *Config
	router *mux.Router
	data   *data
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		router: mux.NewRouter(),
		data:   &data{Result: nil},
	}
}

func (s *APIServer) Start() error {
	log.Println("Server start")
	s.configureRouter()
	return http.ListenAndServe(s.config.bind_addr, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/", s.handleStart)
	s.router.HandleFunc("/calc{num}", s.handleCalc)
	s.router.HandleFunc("/check", s.handleCheck)
}

func (s *APIServer) handleStart(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("task1/templates/index.html")
	if err != nil {
		log.Printf("error to parse template: %+v", err)
		http.Error(w, "server error", 500)
	}

	w.WriteHeader(200)

	err = templ.ExecuteTemplate(w, "index.html", s.data)
	if err != nil {
		log.Printf("error to execute template: %+v", err)
	}
}

func (s *APIServer) handleCalc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("%+v\n", vars)

	num, err := strconv.Atoi(vars["num"])
	if err != nil {
		log.Printf("error to parse num: %+v\n", err)
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	s.data.Result = internal.Calculate(num)

	t, err := template.ParseFiles("task1/templates/index.html")
	if err != nil {
		log.Printf("Error to parse html: %+v", err)
		http.Error(w, "Server error", 500)
		return
	}

	log.Printf("%+v", s.data)

	w.WriteHeader(200)
	t.ExecuteTemplate(w, "index.html", s.data)
	if err != nil {
		log.Printf("template error: %+v", err)
	}
}

func (s *APIServer) handleCheck(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string, 1)
	resp["progress"] = internal.CheckProgress()
	fmt.Println(1)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(2)
	respJson, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error to marshal: %+v\n", err)
	}
	fmt.Println(3)
	w.Write(respJson)
	fmt.Println(w)
	return
}
